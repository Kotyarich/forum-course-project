const baseUrl = 'http://localhost:5000/api/forum';

class ForumService {
  getAll = async () => {
    const url = baseUrl + 's';

    const headers = new Headers();
    const options = {
      method: 'GET',
      headers,
    };

    const request = new Request(url, options);
    const response = await fetch(request);
    return response.json();
  };

  getBySlug = async (slug) => {
    const url = baseUrl + '/' + slug + '/details';

    const headers = new Headers();
    const options = {
      method: 'GET',
      headers,
    };

    const request = new Request(url, options);
    const response = await fetch(request);
    return response.json();
  };
}

export default ForumService;