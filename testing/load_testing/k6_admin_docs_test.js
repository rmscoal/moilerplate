import http from 'k6/http';
import { check, group } from 'k6';

const BASE_URL = __ENV.HOST ? __ENV.HOST : "http://localhost:8082";

export default function() {
  group('01. Login Page', () => {
    check(
      http.get(`${BASE_URL}/api/v1/docs/login`),
      {
        "Status code must be 200": (r) => r.status == 200,
        "Content type must be HTML": (r) => r.headers["Content-Type"].includes("text/html")
      },
    )
  })

  group('02. Post admin secret', () => {
    check(
      http.post(`${BASE_URL}/api/v1/docs/login`, { "key": "verystrongpassword" }),
      {
        "Status code must be 200": (r) => r.status == 200,
      },
    )
  })

  group('03. Accessing index.html', () => {
    check(
      http.get(`${BASE_URL}/api/v1/docs/index.html`),
      {
        "Status code must be 200": (r) => r.status == 200,
      },
    )
  })
}
