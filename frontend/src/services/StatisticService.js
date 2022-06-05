const baseUrl = 'http://localhost:5010/';

class StatisticService {
  get = async () => {
    const url = baseUrl + 'statistic';

    const headers = new Headers();
    const options = {
      method: 'GET',
      headers,
      credentials: 'include',
    };

    const request = new Request(url, options);
    const response = await fetch(request);
    return response.json();
  };
}

export default StatisticService;