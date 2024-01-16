import http from 'k6/http';
import { check, group, fail } from 'k6';

function randomCharacters(length) {
  const charset = 'abcdefghijklmnopqrstuvwxyz';
  let res = '';
  while (length--) res += charset[(Math.random() * charset.length) | 0];
  return res;
}

function randomNumbers(length) {
  const charset = '123456789';
  let res = '';
  while (length--) res += charset[(Math.random() * charset.length) | 0];
  return res;
}

const BASE_URL = __ENV.HOST ? __ENV.HOST : "http://localhost:8082";
const DEFAULT_HEADERS = { "Content-Type": "application/json" };

export function setup() {
  // The data to pass to VU function
  let vuData = {
    username: `${randomCharacters(10)}`,
    firstName: `${randomCharacters(10)}`,
    lastName: `${randomCharacters(10)}`,
    email: `${randomCharacters(10)}@gmail.com`,
    phoneNumber: `+62${randomNumbers(10)}`,
    password: "password",
    accessToken: "",
    refreshToken: "",
  }

  const response = http.post(`${BASE_URL}/api/v1/credentials/signup`, JSON.stringify(vuData), { headers: DEFAULT_HEADERS });
  const signupSuccessful = check(response, { 'Status code must be 201': (r) => r.status == 201 })
  if (signupSuccessful) {
    let data = response.json().data;
    vuData.accessToken = data.accessToken;
    vuData.refreshToken = data.refreshToken;
    return vuData;
  } else {
    fail("Signup wasn't successful")
  }
}

export default function(vuData) {
  group("01. Login", () => {
    const payload = {
      "username": vuData.username,
      "password": vuData.data,
    }

    const headers = Object.assign({ "Authorization": vuData.accessToken }, DEFAULT_HEADERS);
    const response = http.post(`${BASE_URL}/api/v1/credentials/login`, JSON.stringify(payload), { headers });
    check(response, { "Status code must be 200": (r) => r.status == 200 || r.status == 429 });
    if (response.status == 200) {
      let data = response.json().data
      vuData.accessToken = data.accessToken
      vuData.refreshToken = data.refreshToken
    }
  })

  group("02. Profile Management", () => {
    const headers = Object.assign({ "Authorization": vuData.accessToken }, DEFAULT_HEADERS);
    const response = http.get(`${BASE_URL}/api/v1/ptd/profiles/me`, { headers });
    check(response, {
      "Status code must be 200": (r) => r.status == 200,
      "Correct username": (r) => r.status == 200,
    });
  })

  group("03. Change email", () => {
    const payload = {
      "emails": [
        {
          "email": `${randomCharacters(10)}@gmail.com`,
          "isPrimary": true
        },
        {
          "email": `${randomCharacters(10)}@gmail.com`,
          "isPrimary": false
        },
      ]
    }
    const headers = Object.assign({ "Authorization": vuData.accessToken }, DEFAULT_HEADERS);
    const response = http.put(`${BASE_URL}/api/v1/ptd/profiles/email`, JSON.stringify(payload), { headers })
    check(response, {
      "Status code must be 200": (r) => r.status == 200,
    })
  })
}

