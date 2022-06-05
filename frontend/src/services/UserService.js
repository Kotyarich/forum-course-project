const baseUrl = 'http://localhost:5001/user/';
const baseUrlAuth = 'http://localhost:5002/user/';

class UserService {
  signIn = async (user) => {
    const url = baseUrlAuth + 'auth';

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
    const url = baseUrlAuth + user.nickname + '/create';

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
    const url = baseUrlAuth + 'check';
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
      method: 'PATCH',
      headers,
      credentials: 'include',
      body: JSON.stringify({
        about: user.about,
        email: user.email,
        fullname: user.fullname,
      })
    };

    const request = new Request(url, options);
    return await fetch(request);
  };

  singOut = async () => {
    const url = baseUrlAuth + 'signout';
    const options = {method: 'GET', credentials: 'include'};
    const request = new Request(url, options);
    return await fetch(request);
  }
}

export default UserService;