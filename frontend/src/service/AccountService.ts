import http from "../http-common";



const register = (name: string, username: string, email: string, password: string) => {
    return http.post("api/v1/accounts/register", {
        name,
        username,
        email,
        password,
    });
};

const login = (username: string, password: string) => {
    return http.post("api/v1/accounts/login", {
        username,
        password,
    })
        .then((response) => {
            if (response.data.token) {
                localStorage.setItem("user", JSON.stringify(response.data));
            }

            return response.data;
        });
};

const logout = () => {
    localStorage.removeItem("user");
};

const getCurrentUser = () => {
    const userStr = localStorage.getItem("user");
    if (userStr) return JSON.parse(userStr);

    return null;
};

const AccountService = {
    register,
    login,
    logout,
    getCurrentUser,
};

export default AccountService;