import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';
import { BASE_URL, LOAD_STAGES, THRESHOLDS, testUser } from './config.js';

// Custom metrics
const createPostDuration = new Trend('create_post_duration', true);
const likePostDuration = new Trend('like_post_duration', true);
const commentPostDuration = new Trend('comment_post_duration', true);
const getPostDuration = new Trend('get_post_duration', true);
const postsCreated = new Counter('posts_created');
const postErrorRate = new Rate('post_error_rate');

export const options = {
    stages: LOAD_STAGES[__ENV.PROFILE || 'load'],
    thresholds: {
        ...THRESHOLDS,
        create_post_duration: ['p(95)<500', 'p(99)<1500'],
        like_post_duration: ['p(95)<300', 'p(99)<800'],
        comment_post_duration: ['p(95)<400', 'p(99)<1000'],
        get_post_duration: ['p(95)<200', 'p(99)<500'],
        post_error_rate: ['rate<0.05'],
    },
    tags: { testSuite: 'posts' },
};

export default function () {
    const user = testUser('post');
    const params = { headers: { 'Content-Type': 'application/json' } };

    // Signup and login
    http.post(`${BASE_URL}/users/signup`, JSON.stringify(user), params);
    const loginRes = http.post(
        `${BASE_URL}/users/login`,
        JSON.stringify({ username: user.username, password: user.password }),
        params
    );

    if (loginRes.status !== 200) {
        postErrorRate.add(true);
        return;
    }

    sleep(0.3);

    // Create post
    let postId;
    group('Create Post', () => {
        const createRes = http.post(
            `${BASE_URL}/posts`,
            JSON.stringify({
                content_text: `Load test post by ${user.username} at ${Date.now()}. #wandersphere #loadtest`,
            }),
            params
        );

        createPostDuration.add(createRes.timings.duration);
        const ok = check(createRes, {
            'create post status 200': (r) => r.status === 200,
        });
        postErrorRate.add(!ok);

        if (createRes.status === 200) {
            try {
                const body = createRes.json();
                postId = body.post_id || body.id;
                postsCreated.add(1);
            } catch (e) { }
        }
    });

    if (!postId) return;
    sleep(0.3);

    // Get post
    group('Get Post', () => {
        const getRes = http.get(`${BASE_URL}/posts/${postId}`, params);
        getPostDuration.add(getRes.timings.duration);
        check(getRes, {
            'get post status 200': (r) => r.status === 200,
        });
    });

    sleep(0.2);

    // Like post
    group('Like Post', () => {
        const likeRes = http.post(`${BASE_URL}/posts/${postId}/likes`, null, params);
        likePostDuration.add(likeRes.timings.duration);
        check(likeRes, {
            'like post status 200': (r) => r.status === 200,
        });
    });

    sleep(0.2);

    // Comment on post
    group('Comment on Post', () => {
        const commentRes = http.post(
            `${BASE_URL}/posts/${postId}`,
            JSON.stringify({ content: `Load test comment at ${Date.now()}` }),
            params
        );
        commentPostDuration.add(commentRes.timings.duration);
        check(commentRes, {
            'comment post status 200': (r) => r.status === 200,
        });
    });

    sleep(0.5);
}
