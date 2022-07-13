package thousandeyes

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_GetRoles(t *testing.T) {
	setup()
	out := `{"roles": [{"roleName": "admin", "roleId": 2, "hasManagementPermissions": 0, "builtin": 0}, {"roleName": "user1", "roleId": 1, "hasManagementPermissions": 1, "builtin": 1}]}`
	mux.HandleFunc("/roles.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(out))
	})

	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}

	res, err := client.GetRoles()
	if err != nil {
		t.Fatal(err)
	}
	expected := []AccountGroupRole{
		{
			RoleName:                 String("admin"),
			RoleID:                   Int64(2),
			HasManagementPermissions: Bool(false),
			Builtin:                  Bool(false),
		},
		{
			RoleName:                 String("user1"),
			RoleID:                   Int64(1),
			HasManagementPermissions: Bool(true),
			Builtin:                  Bool(true),
		},
	}
	assert.Equal(t, &expected, res)
}

func TestClient_GetRole(t *testing.T) {
	setup()
	out := `{"roles": [{"roleName": "admin", "roleId": 1, "hasManagementPermissions": 0, "builtin": 0}]}`
	mux.HandleFunc("/roles/1.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(out))
	})

	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}

	res, err := client.GetRole(1)
	if err != nil {
		t.Fatal(err)
	}
	expected := AccountGroupRole{
		RoleName:                 String("admin"),
		RoleID:                   Int64(1),
		HasManagementPermissions: Bool(false),
		Builtin:                  Bool(false),
	}
	assert.Equal(t, &expected, res)
}

func TestClient_CreateRole(t *testing.T) {
	setup()
	out := `{"roleName": "ThousandEyes SRE", "roleId": 1000, "hasManagementPermissions": 1, "builtin": 0}`
	mux.HandleFunc("/roles/new.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(out))
	})

	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	create := AccountGroupRole{
		RoleName:                 String("ThousandEyes SRE"),
		HasManagementPermissions: Bool(true),
	}
	res, err := client.CreateRole(create)
	if err != nil {
		t.Fatal(err)
	}

	expected := AccountGroupRole{
		RoleName:                 String("ThousandEyes SRE"),
		HasManagementPermissions: Bool(true),
		Builtin:                  Bool(false),
		RoleID:                   Int64(1000),
	}
	assert.Equal(t, &expected, res)
}

func TestClient_DeleteRole(t *testing.T) {
	setup()
	mux.HandleFunc("/roles/1/delete.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}

	_ = client.DeleteUser(1)
}

func TestClient_UpdateRole(t *testing.T) {
	setup()
	out := `{"roleName": "ThousandEyes SRE", "roleId": 1000, "hasManagementPermissions": 1, "builtin": 0}`
	mux.HandleFunc("/roles/1/update.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(out))
	})

	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	update := AccountGroupRole{
		RoleName:                 String("ThousandEyes SRE"),
		HasManagementPermissions: Bool(true),
	}
	res, err := client.UpdateRole(1, update)
	if err != nil {
		t.Fatal(err)
	}

	expected := AccountGroupRole{
		RoleName:                 String("ThousandEyes SRE"),
		RoleID:                   Int64(1000),
		HasManagementPermissions: Bool(true),
		Builtin:                  Bool(false),
	}
	assert.Equal(t, &expected, res)
}

func TestClient_GetRoleStatusCode(t *testing.T) {
	setup()
	out := `{"roles": [{"roleName": "admin", "roleId": 2, "hasManagementPermissions": 0, "builtin": 0}, {"roleName": "user1", "roleId": 1, "hasManagementPermissions": 1, "builtin": 1}]}`
	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	mux.HandleFunc("/roles.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(out))
	})

	_, err := client.GetRoles()
	teardown()
	assert.ErrorContains(t, err, "Response did not contain formatted error: %!s(<nil>). HTTP response code: 400")
}

func TestClient_CreateRoleStatusCode(t *testing.T) {
	setup()
	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	mux.HandleFunc("/roles/new.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{}`))
	})
	_, err := client.CreateRole(AccountGroupRole{})
	teardown()
	assert.ErrorContains(t, err, "Response did not contain formatted error: %!s(<nil>). HTTP response code: 400")
}

func TestClient_UpdateRoleStatusCode(t *testing.T) {
	setup()
	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	mux.HandleFunc("/roles/1/update.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{}`))
	})
	_, err := client.UpdateRole(1, AccountGroupRole{})
	teardown()
	assert.ErrorContains(t, err, "Response did not contain formatted error: %!s(<nil>). HTTP response code: 400")
}

func TestClient_DeleteRoleStatusCode(t *testing.T) {
	setup()
	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	mux.HandleFunc("/roles/1/delete.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{}`))
	})
	err := client.DeleteRole(1)
	teardown()
	assert.ErrorContains(t, err, "Response did not contain formatted error: %!s(<nil>). HTTP response code: 400")
}

func TestClient_GetRolesJsonError(t *testing.T) {
	out := `{"users": [test]}`
	setup()
	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	mux.HandleFunc("/roles.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		_, _ = w.Write([]byte(out))
	})
	_, err := client.GetRoles()
	assert.Error(t, err)
	assert.EqualError(t, err, "could not decode JSON response: invalid character 'e' in literal true (expecting 'r')")
}

func TestClient_GetRoleJsonError(t *testing.T) {
	out := `{"users": [test]}`
	setup()
	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	mux.HandleFunc("/roles/1.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		_, _ = w.Write([]byte(out))
	})
	_, err := client.GetRole(1)
	assert.Error(t, err)
	assert.EqualError(t, err, "could not decode JSON response: invalid character 'e' in literal true (expecting 'r')")
}

func TestClient_UpdateRolesJsonError(t *testing.T) {
	out := `{"users": [test]}`
	setup()
	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	mux.HandleFunc("/roles/1/update.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		_, _ = w.Write([]byte(out))
	})
	_, err := client.UpdateRole(1, AccountGroupRole{})
	assert.Error(t, err)
	assert.EqualError(t, err, "could not decode JSON response: invalid character 'e' in literal true (expecting 'r')")
}

func TestClient_CreateRoleJsonError(t *testing.T) {
	out := `{"users": [test]}`
	setup()
	var client = &Client{APIEndpoint: server.URL, AuthToken: "foo"}
	mux.HandleFunc("/roles/new.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(out))
	})
	_, err := client.CreateRole(AccountGroupRole{})
	assert.Error(t, err)
	assert.EqualError(t, err, "could not decode JSON response: invalid character 'e' in literal true (expecting 'r')")
}
