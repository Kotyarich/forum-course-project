const baseUrl = 'http://localhost:5000/api/user/';

class UserService {
  signIn = async (user) => {
    const url = baseUrl + 'auth';

    const headers = new Headers();
    headers.append('Content-Type', 'application/json');

    const options = {
      method: 'POST',
      headers,
      credentials: 'include',
      body: JSON.stringify({
        nickname: user.nickname,
        password: user.password,
      })
    };

    const request = new Request(url, options);
    const response = await fetch(request);
    return response.json();
  };

  create = async (user) => {
    const url = baseUrl + user.nickname + '/create';

    const headers = new Headers();
    headers.append('Content-Type', 'application/json');

    const options = {
      method: 'POST',
      headers,
      credentials: 'include',
      body: JSON.stringify({
        about: user.about,
        email: user.email,
        fullname: user.fullname,
        password: user.password,
      })
    };

    const request = new Request(url, options);
    return await fetch(request);
  };

  get = async (nickname) => {
    const url = baseUrl + nickname + '/profile';
    const options = {method: 'GET', credentials: 'include'};

    const request = new Request(url, options);
    const response = await fetch(request);
    return response.json();
  };

  checkAuth = async () => {
    const url = baseUrl + 'check';
    const options = {method: 'GET', credentials: 'include'};

    const request = new Request(url, options);
    const response = await fetch(request);
    return response.json();
  };

  change = async (user) => {
    const url = baseUrl + user.nickname + '/profile';

    const headers = new Headers();
    headers.append('Content-Type', 'application/json');

    const options = {
      method: 'POST',
      headers,
      credentials: 'include',
      body: JSON.stringify({
        about: user.about,
        email: user.email,
        fullname: user.fullName,
      })
    };

    const request = new Request(url, options);
    return await fetch(request);
  };

  singOut = async () => {
    const url = baseUrl + 'signout';
    const options = {method: 'GET', credentials: 'include'};
    const request = new Request(url, options);
    return await fetch(request);
  }
}

export default UserService;