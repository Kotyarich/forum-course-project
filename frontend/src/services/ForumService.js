const baseUrl = 'http://localhost:5000/forum';
const baseThreadUrl = 'http://localhost:5005/forum';

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

  createThread = async (thread, forum_slug) => {
    const url = baseThreadUrl + "/" + forum_slug + "/create";

    const headers = new Headers();
    headers.append('Content-Type', 'application/json');

    const options = {
      method: 'POST',
      headers,
      credentials: 'include',
      body: JSON.stringify({
        author: thread.author,
        created: thread.created,
        message: thread.message,
        title: thread.title,
        slug: forum_slug + '_' + thread.title,
      })
    };

    const request = new Request(url, options);
    return await fetch(request);
  };
}

export default ForumService;