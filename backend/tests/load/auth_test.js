import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { BASE_URL, LOAD_STAGES, THRESHOLDS, testUser } from './config.js';

// Custom metrics
const loginDuration = new Trend('login_duration', true);
const signupDuration = new Trend('signup_duration', true);
const authErrorRate = new Rate('auth_error_rate');

export const options = {
    stages: LOAD_STAGES[__ENV.PROFILE || 'load'],
    thresholds: {
        ...THRESHOLDS,
        login_duration: ['p(95)<300', 'p(99)<800'],
        signup_duration: ['p(95)<500', 'p(99)<1500'],
        auth_error_rate: ['rate<0.05'],
    },
    tags: { testSuite: 'auth' },
};

export default function () {
    const user = testUser('auth');
    const params = { headers: { 'Content-Type': 'application/json' } };

    group('Signup Flow', () => {
        const signupRes = http.post(
            `${BASE_URL}/users/signup`,
            JSON.stringify(user),
            params
        );

        signupDuration.add(signupRes.timings.duration);
        const signupOk = check(signupRes, {
            'signup status 200': (r) => r.status === 200,
            'signup has message': (r) => r.json('message') !== undefined,
        });
        authErrorRate.add(!signupOk);
    });

    sleep(0.5);

    group('Login Flow', () => {
        const loginRes = http.post(
            `${BASE_URL}/users/login`,
            JSON.stringify({
                username: user.username,
                password: user.password,
            }),
            params
        );

        loginDuration.add(loginRes.timings.duration);
        const loginOk = check(loginRes, {
            'login status 200': (r) => r.status === 200,
            'login has session cookie': (r) => {
                const cookies = r.cookies;
                return cookies['wandersphere_session'] !== undefined ||
                    cookies['session_id'] !== undefined;
            },
        });
        authErrorRate.add(!loginOk);
    });

    sleep(0.3);

    group('Get Profile', () => {
        const profileRes = http.get(`${BASE_URL}/users/1`, params);
        check(profileRes, {
            'profile status 200': (r) => r.status === 200,
        });
    });

    sleep(0.5);
}
