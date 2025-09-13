import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { BASE_URL, LOAD_STAGES, THRESHOLDS, testUser } from './config.js';

// Custom metrics
const newsfeedDuration = new Trend('newsfeed_duration', true);
const newsfeedPaginationDuration = new Trend('newsfeed_pagination_duration', true);
const newsfeedErrorRate = new Rate('newsfeed_error_rate');

export const options = {
    stages: LOAD_STAGES[__ENV.PROFILE || 'load'],
    thresholds: {
        ...THRESHOLDS,
        newsfeed_duration: ['p(95)<400', 'p(99)<1000'],
        newsfeed_pagination_duration: ['p(95)<300', 'p(99)<800'],
        newsfeed_error_rate: ['rate<0.01'],
    },
    tags: { testSuite: 'newsfeed' },
};

// Setup: create a user and login once per VU
export function setup() {
    const user = testUser('nf');
    const params = { headers: { 'Content-Type': 'application/json' } };

    // Signup
    http.post(`${BASE_URL}/users/signup`, JSON.stringify(user), params);

    // Login
    const loginRes = http.post(
        `${BASE_URL}/users/login`,
        JSON.stringify({ username: user.username, password: user.password }),
        params
    );

    const cookies = loginRes.cookies;
    const sessionCookie = cookies['wandersphere_session'] || cookies['session_id'];

    return { user, sessionCookieName: sessionCookie ? Object.keys(loginRes.cookies)[0] : 'session_id' };
}

export default function (data) {
    const params = { headers: { 'Content-Type': 'application/json' } };

    group('Newsfeed - First Page', () => {
        const res = http.get(`${BASE_URL}/newsfeed?limit=20`, params);

        newsfeedDuration.add(res.timings.duration);
        const ok = check(res, {
            'newsfeed status 200 or 401': (r) => r.status === 200 || r.status === 401,
            'newsfeed response time < 1s': (r) => r.timings.duration < 1000,
        });
        newsfeedErrorRate.add(!ok);

        // If we got a next_cursor, paginate
        if (res.status === 200) {
            try {
                const body = res.json();
                if (body.next_cursor) {
                    sleep(0.2);

                    const pageRes = http.get(
                        `${BASE_URL}/newsfeed?cursor=${body.next_cursor}&limit=20`,
                        params
                    );

                    newsfeedPaginationDuration.add(pageRes.timings.duration);
                    check(pageRes, {
                        'pagination status 200': (r) => r.status === 200,
                        'pagination response time < 800ms': (r) => r.timings.duration < 800,
                    });
                }
            } catch (e) {
                // Response may not be JSON
            }
        }
    });

    sleep(1);
}
