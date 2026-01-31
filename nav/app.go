package nav

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"wrzapi/frontend"
)

type Config struct {
	DataPath string
	Dev      bool
}

type Category struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Order int32  `json:"order"`
}

type Item struct {
	ID         uint32  `json:"id"`
	Name       string  `json:"name"`
	URL        string  `json:"url"`
	CategoryID *uint32 `json:"category_id"`
	Order      int32   `json:"order"`
	AvatarURL  string  `json:"avatar_url"`
	Summary    string  `json:"summary"`
}

type AdminAuth struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type DataFile struct {
	NextID     uint32     `json:"next_id"`
	Categories []Category `json:"categories"`
	Items      []Item     `json:"items"`
	Admin      AdminAuth  `json:"admin"`
}

type AppState struct {
	mu         sync.Mutex
	dataPath   string
	nextID     uint32
	items      []Item
	categories []Category
	admin      AdminAuth
	sessions   map[string]string
}

type App struct {
	state  *AppState
	mux    *http.ServeMux
	distFS fs.FS
}

func New(cfg Config) (*App, error) {
	dataPath := cfg.DataPath
	if strings.TrimSpace(dataPath) == "" {
		dataPath = "data.json"
	}

	data, err := loadData(dataPath)
	if err != nil {
		data = DataFile{NextID: 1, Categories: []Category{}, Items: []Item{}, Admin: defaultAdmin()}
	}

	state := &AppState{
		dataPath:   dataPath,
		nextID:     data.NextID,
		items:      data.Items,
		categories: data.Categories,
		admin:      data.Admin,
		sessions:   map[string]string{},
	}

	var distFS fs.FS
	if cfg.Dev {
		distFS = os.DirFS("frontend/dist")
	} else {
		sub, err := fs.Sub(frontend.Assets, "dist")
		if err != nil {
			return nil, fmt.Errorf("load nav frontend dist: %w", err)
		}
		distFS = sub
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(distFS))

	mux.Handle("/assets/", fileServer)
	mux.Handle("/favicon.ico", fileServer)
	mux.Handle("/manifest.json", fileServer)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			writeText(w, http.StatusNotFound, "not found")
			return
		}
		serveIndex(w, distFS)
	})

	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin" {
			writeText(w, http.StatusNotFound, "not found")
			return
		}
		cookie, err := r.Cookie("nav_session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		state.mu.Lock()
		_, ok := state.sessions[cookie.Value]
		state.mu.Unlock()
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		serveIndex(w, distFS)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/login" {
			writeText(w, http.StatusNotFound, "not found")
			return
		}
		serveIndex(w, distFS)
	})

	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			state.handleGetData(w)
		case http.MethodPost:
			state.requireAuth(func(w http.ResponseWriter, r *http.Request) {
				state.handleRestore(w, r)
			})(w, r)
		default:
			writeText(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeText(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		state.handleLogin(w, r)
	})

	mux.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeText(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		state.handleLogout(w)
	})

	mux.HandleFunc("/api/password", state.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			writeText(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		state.handleChangePassword(w, r)
	}))

	mux.HandleFunc("/api/category", state.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeText(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		state.handleCreateCategory(w, r)
	}))

	mux.HandleFunc("/api/category/", state.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
		if idStr == "" {
			writeText(w, http.StatusBadRequest, "invalid id")
			return
		}
		idVal, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			writeText(w, http.StatusBadRequest, "invalid id")
			return
		}
		id := uint32(idVal)
		switch r.Method {
		case http.MethodPut:
			state.handleUpdateCategory(w, r, id)
		case http.MethodDelete:
			state.handleDeleteCategory(w, id)
		default:
			writeText(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}))

	mux.HandleFunc("/api/item", state.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeText(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		state.handleCreateItem(w, r)
	}))

	mux.HandleFunc("/api/item/", state.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/item/")
		if idStr == "" {
			writeText(w, http.StatusBadRequest, "invalid id")
			return
		}
		idVal, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			writeText(w, http.StatusBadRequest, "invalid id")
			return
		}
		id := uint32(idVal)
		switch r.Method {
		case http.MethodPut:
			state.handleUpdateItem(w, r, id)
		case http.MethodDelete:
			state.handleDeleteItem(w, id)
		default:
			writeText(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}))

	return &App{state: state, mux: mux, distFS: distFS}, nil
}

func (a *App) Handler() http.Handler {
	return a.mux
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func constantTimeEquals(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func defaultAdmin() AdminAuth {
	return AdminAuth{Username: "admin", PasswordHash: hashPassword("admin")}
}

func loadData(path string) (DataFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DataFile{NextID: 1, Categories: []Category{}, Items: []Item{}, Admin: defaultAdmin()}, nil
		}
		return DataFile{}, err
	}
	if len(raw) == 0 {
		return DataFile{NextID: 1, Categories: []Category{}, Items: []Item{}, Admin: defaultAdmin()}, nil
	}
	var data DataFile
	if err := json.Unmarshal(raw, &data); err == nil {
		if data.NextID == 0 {
			data.NextID = 1
		}
		if data.Admin.Username == "" || data.Admin.PasswordHash == "" {
			data.Admin = defaultAdmin()
		}
		return data, nil
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		return DataFile{NextID: 1, Categories: []Category{}, Items: []Item{}, Admin: defaultAdmin()}, nil
	}

	out := DataFile{NextID: 1, Categories: []Category{}, Items: []Item{}, Admin: defaultAdmin()}
	if v, ok := rawMap["next_id"]; ok {
		var next uint32
		if err := json.Unmarshal(v, &next); err == nil && next > 0 {
			out.NextID = next
		}
	}
	if v, ok := rawMap["categories"]; ok {
		var cats []Category
		if err := json.Unmarshal(v, &cats); err == nil {
			out.Categories = cats
		}
	}
	if v, ok := rawMap["items"]; ok {
		var items []Item
		if err := json.Unmarshal(v, &items); err == nil {
			out.Items = items
		}
	}
	if v, ok := rawMap["admin"]; ok {
		var admin AdminAuth
		if err := json.Unmarshal(v, &admin); err == nil && admin.Username != "" && admin.PasswordHash != "" {
			out.Admin = admin
		}
	}
	return out, nil
}

func (s *AppState) save() error {
	data := DataFile{NextID: s.nextID, Categories: s.categories, Items: s.items, Admin: s.admin}
	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.dataPath, payload, 0644)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func writeText(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(text))
}

func readBody(r *http.Request, limit int64) ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(io.LimitReader(r.Body, limit))
}

func (s *AppState) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("nav_session")
		if err != nil || cookie.Value == "" {
			writeText(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		s.mu.Lock()
		_, ok := s.sessions[cookie.Value]
		s.mu.Unlock()
		if !ok {
			writeText(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next(w, r)
	}
}

func (s *AppState) handleGetData(w http.ResponseWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, DataFile{NextID: s.nextID, Categories: s.categories, Items: s.items, Admin: AdminAuth{Username: s.admin.Username}})
}

func (s *AppState) handleRestore(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r, 5*1024*1024)
	if err != nil {
		writeText(w, http.StatusBadRequest, "invalid body")
		return
	}
	var data DataFile
	if err := json.Unmarshal(body, &data); err != nil {
		writeText(w, http.StatusBadRequest, "invalid json")
		return
	}
	if data.NextID == 0 {
		data.NextID = 1
	}
	if data.Admin.Username == "" || data.Admin.PasswordHash == "" {
		data.Admin = defaultAdmin()
	}

	s.mu.Lock()
	s.nextID = data.NextID
	s.items = data.Items
	s.categories = data.Categories
	s.admin = data.Admin
	err = s.save()
	s.mu.Unlock()
	if err != nil {
		writeText(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *AppState) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r, 512*1024)
	if err != nil {
		writeText(w, http.StatusBadRequest, "invalid body")
		return
	}
	var req Item
	if err := json.Unmarshal(body, &req); err != nil {
		writeText(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.URL) == "" {
		writeText(w, http.StatusBadRequest, "name and url required")
		return
	}

	s.mu.Lock()
	req.ID = s.nextID
	s.nextID++
	s.items = append(s.items, req)
	err = s.save()
	s.mu.Unlock()
	if err != nil {
		writeText(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusCreated, req)
}

func (s *AppState) handleUpdateItem(w http.ResponseWriter, r *http.Request, id uint32) {
	body, err := readBody(r, 512*1024)
	if err != nil {
		writeText(w, http.StatusBadRequest, "invalid body")
		return
	}
	var req Item
	if err := json.Unmarshal(body, &req); err != nil {
		writeText(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.URL) == "" {
		writeText(w, http.StatusBadRequest, "name and url required")
		return
	}

	s.mu.Lock()
	updated := false
	for i := range s.items {
		if s.items[i].ID == id {
			req.ID = id
			s.items[i] = req
			updated = true
			break
		}
	}
	if updated {
		err = s.save()
	}
	s.mu.Unlock()

	if !updated {
		writeText(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		writeText(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusOK, req)
}

func (s *AppState) handleDeleteItem(w http.ResponseWriter, id uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	removed := false
	filtered := s.items[:0]
	for _, item := range s.items {
		if item.ID == id {
			removed = true
			continue
		}
		filtered = append(filtered, item)
	}
	if !removed {
		writeText(w, http.StatusNotFound, "not found")
		return
	}
	s.items = filtered

	if err := s.save(); err != nil {
		writeText(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *AppState) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r, 128*1024)
	if err != nil {
		writeText(w, http.StatusBadRequest, "invalid body")
		return
	}
	var req Category
	if err := json.Unmarshal(body, &req); err != nil {
		writeText(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeText(w, http.StatusBadRequest, "name required")
		return
	}

	s.mu.Lock()
	req.ID = s.nextID
	s.nextID++
	s.categories = append(s.categories, req)
	err = s.save()
	s.mu.Unlock()
	if err != nil {
		writeText(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusCreated, req)
}

func (s *AppState) handleUpdateCategory(w http.ResponseWriter, r *http.Request, id uint32) {
	body, err := readBody(r, 128*1024)
	if err != nil {
		writeText(w, http.StatusBadRequest, "invalid body")
		return
	}
	var req Category
	if err := json.Unmarshal(body, &req); err != nil {
		writeText(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeText(w, http.StatusBadRequest, "name required")
		return
	}

	s.mu.Lock()
	updated := false
	for i := range s.categories {
		if s.categories[i].ID == id {
			req.ID = id
			s.categories[i] = req
			updated = true
			break
		}
	}
	if updated {
		err = s.save()
	}
	s.mu.Unlock()

	if !updated {
		writeText(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		writeText(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusOK, req)
}

func (s *AppState) handleDeleteCategory(w http.ResponseWriter, id uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := s.categories[:0]
	removed := false
	for _, cat := range s.categories {
		if cat.ID == id {
			removed = true
			continue
		}
		filtered = append(filtered, cat)
	}
	if !removed {
		writeText(w, http.StatusNotFound, "not found")
		return
	}
	s.categories = filtered

	for i := range s.items {
		if s.items[i].CategoryID != nil && *s.items[i].CategoryID == id {
			s.items[i].CategoryID = nil
		}
	}

	if err := s.save(); err != nil {
		writeText(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *AppState) handleLogin(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r, 64*1024)
	if err != nil {
		writeText(w, http.StatusBadRequest, "invalid body")
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeText(w, http.StatusBadRequest, "invalid json")
		return
	}

	s.mu.Lock()
	admin := s.admin
	s.mu.Unlock()

	if req.Username != admin.Username || !constantTimeEquals(hashPassword(req.Password), admin.PasswordHash) {
		writeText(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		writeText(w, http.StatusInternalServerError, "token error")
		return
	}
	session := hex.EncodeToString(token)

	s.mu.Lock()
	s.sessions[session] = admin.Username
	s.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "nav_session",
		Value:    session,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *AppState) handleLogout(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "nav_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *AppState) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r, 64*1024)
	if err != nil {
		writeText(w, http.StatusBadRequest, "invalid body")
		return
	}
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeText(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.NewPassword) == "" {
		writeText(w, http.StatusBadRequest, "new password required")
		return
	}

	s.mu.Lock()
	if !constantTimeEquals(hashPassword(req.OldPassword), s.admin.PasswordHash) {
		s.mu.Unlock()
		writeText(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	s.admin.PasswordHash = hashPassword(req.NewPassword)
	err = s.save()
	s.mu.Unlock()
	if err != nil {
		writeText(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func serveAsset(w http.ResponseWriter, dist fs.FS, path string, contentType string) {
	data, err := fs.ReadFile(dist, path)
	if err != nil {
		writeText(w, http.StatusNotFound, "not found")
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func serveIndex(w http.ResponseWriter, dist fs.FS) {
	serveAsset(w, dist, "index.html", "text/html; charset=utf-8")
}
