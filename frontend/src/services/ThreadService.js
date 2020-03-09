const baseUrl = 'http://localhost:5000/api/';

class ThreadService {
  getAll = async (forumSlug, since = '', limit = 10, offset = 0, desc = true) => {
    console.log("getAll: ", limit, offset);
    const url = baseUrl + 'forum/' + forumSlug + '/threads';
    const urlWithParams = url + '?desc=' + desc
      + ';limit=' + limit + ';offset=' + offset;

    const headers = new Headers();
    const options = {
      method: 'GET',
      headers,
    };

    const request = new Request(urlWithParams, options);
    const response = await fetch(request);
    return response.json();
  };

  vote = async (threadSlug, nickname, vote) => {
    const url = baseUrl + 'thread/' + threadSlug + '/vote';

    const headers = new Headers();
    const options = {
      method: 'POST',
      headers,
      credentials: 'include',
      body: JSON.stringify({
        voice: vote,
        nickname: nickname,
      }),
    };

    const request = new Request(url, options);
    const response = await fetch(request);
    return response.json();
  };
}

export default ThreadService;