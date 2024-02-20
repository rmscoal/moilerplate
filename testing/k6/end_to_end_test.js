// More examples can be seen here:
// https://grafana.com/docs/k6/latest/examples/functional-testing/

import chai, { describe, expect } from "https://jslib.k6.io/k6chaijs/4.3.4.3/index.js";
import { Httpx } from "https://jslib.k6.io/httpx/0.0.6/index.js";
import { generateNewUser } from "./helpers.js";

// Log when failures occurs
chai.config.logFailures = true;

export let options = {
  thresholds: {
    // fail the test if any checks fail or any requests fail
    checks: ["rate == 1.00"],
    http_req_failed: ["rate == 0.00"],
  },
  vus: 1,
  iterations: 1,
};

// Session that makes the HTTP calls for this integration test
const session = new Httpx({ baseURL: __ENV.HOST ? __ENV.HOST : "http://localhost:8082/api/v1" });
// Store test data in a hash map
const hashmap = new Map();

function userCredentials() {
  describe("I am able to signup as new user", () => {
    const user = generateNewUser();
    const payload = JSON.stringify(user);
    const headers = { "Content-Type": "application/json" };
    const response = session.post("/credentials/signup", payload, { headers });

    const json = response.json();
    expect(response.status, "status code should be 201").to.be.equal(201);
    expect(json.status).to.be.equal("OK");
    expect(json.data.username).to.be.equal(user.username);
    hashmap.set('username', user.username);
    hashmap.set('password', user.password);
  });

  describe("I am able to login after signing up", () => {
    const payload = JSON.stringify({
      username: hashmap.get('username'),
      password: hashmap.get('password'),
    });
    const headers = { "Content-Type": "application/json" };
    const response = session.post("/credentials/login", payload, { headers });

    const json = response.json();
    expect(response.status, "status code should be 200").to.be.equal(200);
    expect(json.data.accessToken, "access token should not be empty").to.be.not.empty;
    expect(json.data.refreshToken, "refresh token should not be empty").to.be.not.empty;
    hashmap.set('accessToken', json.data.accessToken);
    hashmap.set('refreshToken', json.data.refreshToken);
  });

  describe("I am able to refresh my access using refresh token", () => {
    const payload = JSON.stringify({
      refreshToken: hashmap.get('refreshToken'),
    });
    const headers = { "Content-Type": "application/json" };
    const response = session.post("/credentials/refresh", payload, { headers });

    const json = response.json();
    expect(response.status, "status code should be 200").to.be.equal(200);
    expect(json.data.accessToken, "access token should not be empty").to.be.not.empty;
    expect(json.data.refreshToken, "refresh token should not be empty").to.be.not.empty;
    hashmap.set('accessToken', json.data.accessToken);
    hashmap.set('refreshToken', json.data.refreshToken);
  });
}

export default function testSuite() {
  userCredentials();
}
