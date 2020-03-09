package http

import (
	"context"
	"dbProject/common"
	"dbProject/models"
	"dbProject/user"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Handler struct {
	useCase user.UseCase
}

func NewHandler(useCase user.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type userInput struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
	Password string `json:"password,omitempty"`
}

func userInputToModel(user userInput) *models.User {
	return &models.User{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		Password: user.Password,
		Email:    user.Email,
		About:    user.About,
	}
}

type signInInput struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func (h *Handler) SignOutHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	cookie, err := request.Cookie("Auth")
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.SignOut(request.Context(), cookie.Value)
	http.SetCookie(writer, &http.Cookie{
		Name:     "Auth",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) UserAuthHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	// read body
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	// parse body
	var input signInInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	ctx := context.WithValue(request.Context(), "UserAgent", request.UserAgent())

	u, token, err := h.useCase.SignIn(ctx, input.Nickname, input.Password)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusForbidden)
		return
	}

	cookie := &http.Cookie{
		Name:     "Auth",
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Path:     "/",
		Domain:   "localhost",
	}
	http.SetCookie(writer, cookie)

	data, err := json.Marshal(UserToUserOutput(u))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	common.WriteData(writer, http.StatusOK, data)
}

func (h *Handler) UserCheckAuthHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	u := request.Context().Value("user")
	if u == nil {
		msg, _ := json.Marshal(map[string]string{"error": "not authorised"})
		common.WriteData(writer, http.StatusForbidden, msg)
		return
	}

	data, err := json.Marshal(UserToUserOutput(u.(*models.User)))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	common.WriteData(writer, http.StatusOK, data)
}

func (h *Handler) UserGetHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	nickname := ps["nickname"]
	u, err := h.useCase.GetProfile(request.Context(), nickname)
	if err != nil {
		if err == user.ErrUserNotFound {
			msg, _ := json.Marshal(map[string]string{"message": "404"})
			common.WriteData(writer, http.StatusNotFound, msg)
			return
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data, err := json.Marshal(UserToUserOutput(u))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	common.WriteData(writer, http.StatusOK, data)
}

func (h *Handler) UserPostHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	nickname := ps["nickname"]

	// read body
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	// parse body
	var input userInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	input.Nickname = nickname

	newUser, err := h.useCase.ChangeProfile(request.Context(), userInputToModel(input))

	if err == user.ErrUserAlreadyExists {
		msg, _ := json.Marshal(map[string]string{"message": "conflict"})
		common.WriteData(writer, http.StatusConflict, msg)
		return
	} else if err == user.ErrUserNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(UserToUserOutput(newUser))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	common.WriteData(writer, http.StatusOK, data)
	return
}

func (h *Handler) UserCreateHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	if request.Method == http.MethodOptions {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		return
	}
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	u := userInput{}
	if err = json.Unmarshal(body, &u); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	model := userInputToModel(u)
	model.Nickname = ps["nickname"]
	ctx := context.WithValue(request.Context(), "UserAgent", request.UserAgent())
	conflicts, token, err := h.useCase.SignUp(ctx, model)
	if err != nil {
		if conflicts == nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := UsersToJsonArray(conflicts)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusConflict)
		_, _ = writer.Write(result)
		return
	}

	cookie := &http.Cookie{
		Name:     "Auth",
		Value:    token,
		HttpOnly: true,
	}
	http.SetCookie(writer, cookie)

	data, err := json.Marshal(UserToUserOutput(model))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write(data)
}

type UserOutput struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
}

func UserToUserOutput(user *models.User) *UserOutput {
	return &UserOutput{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		Email:    user.Email,
		About:    user.About,
	}
}

func UsersToJsonArray(users []*models.User) ([]byte, error) {
	result := []byte{'['}
	for i := 0; i < len(users); i++ {
		if len(result) > 1 {
			result = append(result, ',')
		}

		userOutput := UserToUserOutput(users[i])
		data, err := json.Marshal(userOutput)
		if err != nil {
			return nil, err
		}

		result = append(result, data...)
	}
	result = append(result, ']')

	return result, nil
}
