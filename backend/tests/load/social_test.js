import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { BASE_URL, LOAD_STAGES, THRESHOLDS, testUser } from './config.js';

// Custom metrics
const followDuration = new Trend('follow_duration', true);
const unfollowDuration = new Trend('unfollow_duration', true);
const getFollowersDuration = new Trend('get_followers_duration', true);
const socialErrorRate = new Rate('social_error_rate');

export const options = {
    stages: LOAD_STAGES[__ENV.PROFILE || 'smoke'],
    thresholds: {
        ...THRESHOLDS,
        follow_duration: ['p(95)<400', 'p(99)<1000'],
        unfollow_duration: ['p(95)<400', 'p(99)<1000'],
        get_followers_duration: ['p(95)<300', 'p(99)<800'],
        social_error_rate: ['rate<0.05'],
    },
    tags: { testSuite: 'social' },
};

export default function () {
    const user1 = testUser('social1');
    const user2 = testUser('social2');
    const params = { headers: { 'Content-Type': 'application/json' } };

    // Create two users
    http.post(`${BASE_URL}/users/signup`, JSON.stringify(user1), params);
    const signup2 = http.post(`${BASE_URL}/users/signup`, JSON.stringify(user2), params);

    let user2Id;
    try {
        user2Id = signup2.json('user_id') || signup2.json('id') || 2;
    } catch (e) {
        user2Id = 2;
    }

    // Login as user1
    http.post(
        `${BASE_URL}/users/login`,
        JSON.stringify({ username: user1.username, password: user1.password }),
        params
    );

    sleep(0.3);

    group('Follow User', () => {
        const followRes = http.post(`${BASE_URL}/friends/${user2Id}`, null, params);
        followDuration.add(followRes.timings.duration);
        const ok = check(followRes, {
            'follow status 200': (r) => r.status === 200,
        });
        socialErrorRate.add(!ok);
    });

    sleep(0.3);

    group('Get Followers', () => {
        const followersRes = http.get(`${BASE_URL}/friends/${user2Id}/followers`, params);
        getFollowersDuration.add(followersRes.timings.duration);
        check(followersRes, {
            'get followers status 200': (r) => r.status === 200,
        });
    });

    sleep(0.3);

    group('Get Following', () => {
        const followingRes = http.get(`${BASE_URL}/friends/${user2Id}/followings`, params);
        check(followingRes, {
            'get following status 200': (r) => r.status === 200,
        });
    });

    sleep(0.3);

    group('Unfollow User', () => {
        const unfollowRes = http.del(`${BASE_URL}/friends/${user2Id}`, null, params);
        unfollowDuration.add(unfollowRes.timings.duration);
        check(unfollowRes, {
            'unfollow status 200': (r) => r.status === 200,
        });
    });

    sleep(0.5);
}
