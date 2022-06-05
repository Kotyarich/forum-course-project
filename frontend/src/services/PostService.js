const baseUrl = 'http://localhost:5003/';

class PostService {
  getByThreadSlug = async (slug, limit = 10, sort = "flat", desc = false,
                           offset = 0, since = 0) => {
    const url = baseUrl + 'thread/' + slug + '/posts';
    const urlWithParams = url + '?desc=' + desc + ';offset=' + offset
      + ';limit=' + limit + ';since=' + since + ';sort=' + sort;

    const headers = new Headers();
    const options = {
      method: 'GET',
      headers,
    };

    const request = new Request(urlWithParams, options);
    const response = await fetch(request);
    return response.json();
  };

  create = async (slug, author, message, parent) => {
    const url = baseUrl + 'thread/' + slug + '/create';
    const headers = new Headers();
    const options = {
      method: 'POST',
      headers,
      body: JSON.stringify(
        [{
          author: author,
          parent: +parent,
          message: message
        }]
      ),
    };
    console.log(options);

    const request = new Request(url, options);
    const response = await fetch(request);
    console.log(response);
    return response.json();
  };

  change = async (id, message) => {
    const url = baseUrl + 'post/' + id + '/details';
    const headers = new Headers();
    headers.append("content-type", 'application/json');
    const options = {
      method: 'PUT',
      headers,
      body: JSON.stringify({message: message}),
    };

    const request = new Request(url, options);
    const response = await fetch(request);
    console.log(response);
    return response.json();
  };
}

export default PostService;