import http from 'k6/http';
import { check, sleep } from 'k6';

// Configuration
export const options = {
  vus: 100,           // 10 virtual users
  duration: '10s',   // Run for 30 seconds
};

// Get endpoint from environment variable or use default
const BASE_URL = __ENV.BASE_URL || 'http://garmin-service:5051';
const ENDPOINT = __ENV.ENDPOINT || '/health-check';

export default function () {
  const response = http.get(`${BASE_URL}${ENDPOINT}`);
  
  // Validate status code = 200
  check(response, {
    'status is 200': (r) => {
      if (r.status === 200) {
        return true;
      }
      console.log(`❌ Request failed: status=${r.status}, url=${BASE_URL}${ENDPOINT}, error=${r.error}`);
      return false;
    },
  });
//   sleep(1);
}
