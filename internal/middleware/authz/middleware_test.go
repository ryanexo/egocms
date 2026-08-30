package authz

import (
    "context"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type testUser struct {
    id   uint64
    role string
}

func (user testUser) UserID() uint64 { return user.id }
func (user testUser) Role() string   { return user.role }

type tokenParserFunc func(context.Context, string) (User, error)

func (function tokenParserFunc) Parse(ctx context.Context, token string) (User, error) {
    return function(ctx, token)
}

type permissionCheckerFunc func(context.Context, string, string, string) (bool, error)

func (function permissionCheckerFunc) Check(ctx context.Context, subject, object, action string) (bool, error) {
    return function(ctx, subject, object, action)
}

func TestMiddlewarePolicies(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    parser := tokenParserFunc(func(_ context.Context, token string) (User, error) {
        if token != "valid" {
            return nil, errors.New("invalid token")
        }
        return testUser{id: 7, role: "editor"}, nil
    })
    checker := permissionCheckerFunc(func(_ context.Context, subject, object, action string) (bool, error) {
        return subject == "editor" && object == "article" && action == "create", nil
    })
    auth := New(parser, checker).Resource("article")
    
    tests := []struct {
        name    string
        path    string
        token   string
        wrap    func(gin.HandlerFunc) gin.HandlerFunc
        status  int
        invoked bool
    }{
        {name: "public", path: "/view", wrap: auth.Public().Wrap, status: http.StatusNoContent, invoked: true},
        {name: "credential missing", path: "/draft", wrap: auth.Wrap, status: http.StatusUnauthorized},
        {name: "credential invalid", path: "/draft", token: "invalid", wrap: auth.Wrap, status: http.StatusUnauthorized},
        {name: "credential valid", path: "/draft", token: "valid", wrap: auth.Wrap, status: http.StatusNoContent, invoked: true},
        {name: "permission allowed", path: "/create", token: "valid", wrap: auth.Permission("create").Wrap, status: http.StatusNoContent, invoked: true},
        {name: "permission denied", path: "/delete", token: "valid", wrap: auth.Permission("delete").Wrap, status: http.StatusForbidden},
    }
    
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            invoked := false
            engine := testEngine()
            engine.POST(test.path, test.wrap(func(ctx *gin.Context) {
                invoked = true
                if test.token == "valid" && test.path == "/draft" {
                    user, err := GetCurrentUser(ctx)
                    require.NoError(t, err)
                    requestUser, err := GetCurrentUser(ctx.Request.Context())
                    require.NoError(t, err)
                    require.Equal(t, user.UserID(), requestUser.UserID())
                }
                ctx.Status(http.StatusNoContent)
            }))
            
            request := httptest.NewRequest(http.MethodPost, test.path, nil)
            if test.token != "" {
                request.Header.Set("Authorization", "Bearer "+test.token)
            }
            response := httptest.NewRecorder()
            engine.ServeHTTP(response, request)
            
            assert.Equal(t, test.status, response.Code)
            assert.Equal(t, test.invoked, invoked)
        })
    }
}

func TestPublicDoesNotUseDependencies(t *testing.T) {
    gin.SetMode(gin.TestMode)
    engine := gin.New()
    auth := New(nil, nil).Resource("article")
    engine.GET("/view", auth.Public().Wrap(func(ctx *gin.Context) {
        ctx.Status(http.StatusNoContent)
    }))
    
    request := httptest.NewRequest(http.MethodGet, "/view", nil)
    response := httptest.NewRecorder()
    engine.ServeHTTP(response, request)
    
    assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestPolicyDerivationDoesNotMutateResource(t *testing.T) {
    auth := New(nil, nil).Resource("article")
    
    assert.False(t, auth.public)
    assert.False(t, auth.permission)
    assert.True(t, auth.Public().public)
    assert.Equal(t, "create", auth.Permission("create").action)
    assert.False(t, auth.public)
    assert.Empty(t, auth.action)
}

func testEngine() *gin.Engine {
    engine := gin.New()
    engine.Use(func(ctx *gin.Context) {
        ctx.Next()
        lastError := ctx.Errors.Last()
        if lastError == nil {
            return
        }
        if errors.Is(lastError, ErrAuthorized) {
            ctx.Status(http.StatusUnauthorized)
        }
        if errors.Is(lastError, ErrAccessDenied) {
            ctx.Status(http.StatusForbidden)
        }
    })
    return engine
}
