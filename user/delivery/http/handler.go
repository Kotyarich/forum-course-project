package http

import (
	"dbProject/models"
	"dbProject/user"
	"dbProject/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
		Email:    user.Email,
		About:    user.About,
	}
}

type signInInput struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func (h *Handler) UserAuthHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	// TODO implement
}

func (h *Handler) UserGetHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	nickname := ps["nickname"]
	u, err := h.useCase.GetProfile(request.Context(), nickname)
	if err != nil {
		if err == user.ErrUserNotFound {
			msg, _ := json.Marshal(map[string]string{"message": "404"})
			utils.WriteData(writer, http.StatusNotFound, msg)
			return
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data, err := json.Marshal(userToUserOutput(u))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteData(writer, http.StatusOK, data)
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
		utils.WriteData(writer, http.StatusConflict, msg)
		return
	} else if err == user.ErrUserNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userToUserOutput(newUser))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteData(writer, http.StatusOK, data)
	return
}

func (h *Handler) UserCreateHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
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

	if conflicts, err := h.useCase.SignUp(request.Context(), model); err != nil {
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

	data, err := json.Marshal(userToUserOutput(model))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write(data)
}

type userOutput struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
}

func userToUserOutput(user *models.User) *userOutput {
	return &userOutput{
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

		userOutput := userToUserOutput(users[i])
		data, err := json.Marshal(userOutput)
		if err != nil {
			return nil, err
		}

		result = append(result, data...)
	}
	result = append(result, ']')

	return result, nil
}
