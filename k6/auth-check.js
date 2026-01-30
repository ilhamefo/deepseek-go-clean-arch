import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

// Custom metrics
const loginSuccessRate = new Rate('login_success_rate');
const loginDuration = new Trend('login_duration');
const loginErrors = new Counter('login_errors');

// Test configuration
export const options = {
  // Scenario 1: Smoke test (verifikasi basic functionality)
  // stages: [
  //   { duration: '1m', target: 10 }, // Ramp up to 10 users
  //   { duration: '2m', target: 10 }, // Stay at 10 users
  //   { duration: '1m', target: 0 },  // Ramp down
  // ],

  // Scenario 2: Load test (test normal load)
  // stages: [
  //   { duration: '2m', target: 100 },  // Ramp up to 100 users
  //   { duration: '5m', target: 100 },  // Stay at 100 users
  //   { duration: '2m', target: 0 },    // Ramp down
  // ],

  // Scenario 3: Stress test (find breaking point)
  stages: [
    // { duration: '5m', target: 500 },   // Ramp up to 200 users
    { duration: '5s', target: 100 },   // Stay at 200 users
    // { duration: '2m', target: 400 },   // Spike to 400 users
    // { duration: '5m', target: 400 },   // Stay at 400 users
    // { duration: '2m', target: 0 },     // Ramp down
  ],

  // Scenario 4: Spike test (sudden traffic surge)
  // stages: [
  //   { duration: '10s', target: 100 }, // Normal load
  //   { duration: '1m', target: 1000 }, // Sudden spike
  //   { duration: '3m', target: 1000 }, // Stay at spike
  //   { duration: '10s', target: 100 }, // Return to normal
  //   { duration: '1m', target: 0 },    // Ramp down
  // ],

  thresholds: {
    // HTTP request duration should be below 500ms for 95% of requests
    http_req_duration: ['p(95)<500', 'p(99)<1000'],

    // 99% of requests should succeed
    login_success_rate: ['rate>0.99'],

    // Less than 1% error rate
    http_req_failed: ['rate<0.01'],

    // Custom login duration threshold
    login_duration: ['p(95)<400'],
  },
};

// Base URL
const BASE_URL = 'http://host.docker.internal:5052';

// Test data
const credentials = {
  email: 'ilham@oninyon.com',
  password: 'password',
};

// Main test function
export default function () {
  const url = `${BASE_URL}/auth/login`;

  const payload = JSON.stringify(credentials);

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
    tags: {
      name: 'LoginAPI',
    },
  };

  // Measure login duration
  const startTime = new Date();
  const response = http.post(url, payload, params);
  const endTime = new Date();
  const duration = endTime - startTime;

  // Record custom metrics
  loginDuration.add(duration);

  // Validate response
  const checkResult = check(response, {
    'status is 200': (r) => r.status === 200,
    'response has data field': (r) => r.json().hasOwnProperty('data'),
    'response has message field': (r) => r.json().hasOwnProperty('message'),
    'response has status field': (r) => r.json().hasOwnProperty('status'),
    'message is success': (r) => r.json('message') === 'success',
    'status code in body is 200': (r) => r.json('status') === 200,
    // 'response time < 500ms': (r) => r.timings.duration < 500,
    // 'response time < 1000ms': (r) => r.timings.duration < 1000,
  });

  // Record success rate
  loginSuccessRate.add(checkResult);

  // Count errors
  if (!checkResult || response.status !== 200) {
    loginErrors.add(1);
    console.error(`Login failed: ${response.status} - ${response.body}`);
  }

  // Think time (simulate real user behavior)
  sleep(1); // 1 second between requests
}

// Setup function (runs once before test)
export function setup() {
  console.log('Starting load test...');
  console.log(`Target: ${BASE_URL}/auth/login`);
  console.log(`Credentials: ${credentials.email}`);

  // Optional: Verify endpoint is reachable
  const healthCheck = http.get(`${BASE_URL}/swagger/index.html`);
  if (healthCheck.status !== 200) {
    console.warn('Warning: Server might not be ready');
  }

  return { startTime: new Date() };
}

// Teardown function (runs once after test)
export function teardown(data) {
  const endTime = new Date();
  const duration = (endTime - data.startTime) / 1000;
  console.log(`Test completed in ${duration.toFixed(2)} seconds`);
}