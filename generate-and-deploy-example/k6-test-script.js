import http from "k6/http";
import { sleep } from "k6";

export let options = {
  scenarios: {
    ramp_data_endpoint: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: '1m', target: 4 },    // Ramp up to 4 users over 1 minute
        { duration: '1m', target: 4 },    // Stay at 4 users for 1 minute
        { duration: '2m', target: 8 },    // Ramp up to 8 users over 2 minutes
        { duration: '2m', target: 15 },   // Ramp up to 15 users over 2 minutes
        { duration: '7m', target: 15 },   // Stay at 15 users for 7 minutes
        { duration: '2m', target: 0 },    // Ramp down to 0 users over 2 minutes
      ],
      exec: "dataEndpoint", // Function to execute
    },
    ramp_health_endpoint: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
      { duration: '2m', target: 2 },    // Ramp up to 2 users over 2 minutes
      { duration: '30s', target: 2 },   // Stay at 2 users for 30 seconds
      { duration: '1m30s', target: 5 }, // Ramp up to 5 users over 1 minute 30 seconds
      { duration: '1m', target: 10 },   // Ramp up to 10 users over 1 minute
      { duration: '9m', target: 10 },   // Stay at 10 users for 9 minutes
      { duration: '1m', target: 0 },    // Ramp down to 0 users over 1 minute
      ],
      exec: "healthEndpoint", // Function to execute
    },
  },
};

export function dataEndpoint() {
  http.get("http://web-server:5001/data");
  sleep(1);
}

export function healthEndpoint() {
  http.get("http://web-server:5001/health");
  sleep(1);
}
