import http from 'k6/http';
import { check, sleep } from 'k6';

// Test configuration
export const options = {
    stages: [
        { duration: '10s', target: 1000 },  // Ramp up to 10 users over 30s
        // { duration: '1m', target: 10 },   // Stay at 10 users for 1m
        // { duration: '30s', target: 0 },   // Ramp down to 0 users
    ],
    // thresholds: {
    //     http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
    //     http_req_failed: ['rate<0.01'],   // Error rate must be less than 1%
    // },
};

// Base URL - adjust based on your server configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:5051';

export default function () {
    // Test health check endpoint
    const res = http.get(`${BASE_URL}/activity/21496702565`);

    // Verify response
    check(res, {
        'status is 200': (r) => {
            if (r.status === 200) {
                return true;
            } else {
                console.error(`Health check failed with status: ${r.status}`);
                console.error(`Health check failed with body: ${r.body}`);
                return false;
            }
        },
        // 'response time < 500ms': (r) => r.timings.duration < 500,
    });

    // sleep(1); // Wait 1 second between iterations
}
