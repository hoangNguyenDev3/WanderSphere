// Shared configuration for all load tests
export const BASE_URL = __ENV.BASE_URL || 'http://localhost:19003/api/v1';

// Load test stages - ramp up, sustain, ramp down
export const LOAD_STAGES = {
    smoke: [
        { duration: '30s', target: 5 },
        { duration: '1m', target: 5 },
        { duration: '30s', target: 0 },
    ],
    load: [
        { duration: '1m', target: 50 },
        { duration: '3m', target: 50 },
        { duration: '1m', target: 100 },
        { duration: '3m', target: 100 },
        { duration: '2m', target: 0 },
    ],
    stress: [
        { duration: '1m', target: 100 },
        { duration: '3m', target: 200 },
        { duration: '3m', target: 500 },
        { duration: '5m', target: 500 },
        { duration: '2m', target: 0 },
    ],
    spike: [
        { duration: '30s', target: 10 },
        { duration: '10s', target: 500 },
        { duration: '1m', target: 500 },
        { duration: '30s', target: 10 },
        { duration: '1m', target: 0 },
    ],
};

// Performance thresholds
export const THRESHOLDS = {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],  // 95th < 500ms, 99th < 1s
    http_req_failed: ['rate<0.01'],                     // Error rate < 1%
    http_reqs: ['rate>100'],                             // Throughput > 100 RPS
};

// Helper to generate unique username
export function uniqueUsername(prefix) {
    return `${prefix}_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`;
}

// Helper to generate test user data
export function testUser(prefix) {
    const username = uniqueUsername(prefix || 'loadtest');
    return {
        username: username,
        password: 'LoadTest123!',
        firstname: 'Load',
        lastname: 'Tester',
        dateofbirth: '1990-01-15',
        email: `${username}@loadtest.com`,
    };
}
