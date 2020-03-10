const baseUrl = 'http://localhost:5000/api/';

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
}

export default PostService;