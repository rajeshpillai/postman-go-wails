// main.js

async function invoke(method, params) {
    return window.backend[method](params);
  }
  
  const loginContainer = document.getElementById('login-container');
  const appContainer = document.getElementById('app-container');
  const loginForm = document.getElementById('login-form');
  const usernameInput = document.getElementById('username');
  const passwordInput = document.getElementById('password');
  const logoutButton = document.getElementById('logout-button');
  const welcomeMessage = document.getElementById('welcome-message');
  
  async function checkLogin() {
    try {
      const username = await window.go.main.App.CheckLogin();
      console.log('Check login result:', username);
      if (username) {
        loginContainer.classList.add('hidden');
        appContainer.classList.remove('hidden');
        welcomeMessage.textContent = `Welcome, ${username}`;
      } else {
        loginContainer.classList.remove('hidden');
        appContainer.classList.add('hidden');
      }
    } catch (error) {
      console.error('Error checking login:', error);
    }
  }
  
  loginForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    const username = usernameInput.value;
    const password = passwordInput.value;
  
    try {
      await window.go.main.App.Login(username, password);
      console.log('Login successful');
      checkLogin();
    } catch (error) {
      console.error('Login error:', error);
      alert('Incorrect username or password');
    }
  });
  
  logoutButton.addEventListener('click', async () => {
    try {
      await window.go.main.App.Logout();
      console.log('Logout successful');
      checkLogin();
    } catch (error) {
      console.error('Logout error:', error);
    }
  });
  
  // Initialize the login check
  checkLogin();
  
  const requestForm = document.getElementById('request-form');
  const urlInput = document.getElementById('url');
  const methodSelect = document.getElementById('method');
  const headersTextArea = document.getElementById('headers');
  const requestBodyTextArea = document.getElementById('request-body');
  const responseTextArea = document.getElementById('response');
  
  requestForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    const url = urlInput.value;
    const method = methodSelect.value;
  
    let headers = {};
    try {
      headers = JSON.parse(headersTextArea.value || '{}');
    } catch (e) {
      responseTextArea.value = 'Invalid JSON in headers';
      console.error('Invalid JSON in headers:', headersTextArea.value);
      return;
    }
  
    const requestBody = requestBodyTextArea.value;
  
    console.log('URL:', url);
    console.log('Method:', method);
    console.log('Headers:', headers);
    console.log('Request Body:', requestBody);
  
    try {
      const response = await window.go.main.App.PerformFetch({
        url,
        method,
        headers,
        body: method !== 'GET' && method !== 'HEAD' ? requestBody : null
      });
  
      const formattedResponse = `
  Status: ${response.status}
  Headers: ${JSON.stringify(response.headers, null, 2)}
  Body: ${response.body}
  `;
      responseTextArea.value = formattedResponse;
    } catch (error) {
      responseTextArea.value = `Error: ${error}`;
      console.error('Fetch error:', error);
    }
  });
  