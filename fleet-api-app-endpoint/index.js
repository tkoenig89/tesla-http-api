export default {
    async fetch(request, env, ctx) {
        async function gatherResponse(response) {
            const { headers } = response;
            const contentType = headers.get('content-type') || '';
            if (contentType.includes('application/json')) {
                return await response.json();
            }
            return response.text();
        }

        const url = new URL(request.url);
        switch (url.pathname) {
            case '/.well-known/appspecific/com.tesla.3p.public-key.pem':
                return new Response(env.PUBLIC_KEY);
            case '/setup':
                return Response.redirect(`https://www.tesla.com/_ak/${url.hostname}`, 301);
            case '/revoke':
                return Response.redirect(`https://auth.tesla.com/user/revoke/consent?revoke_client_id=${env.CLIENT_ID}&back_url=https://accounts.tesla.com/account-settings/security?tab=tpty-apps`, 301);
            case '/callback':
                const code = url.searchParams.get('code');
                if (!code) {
                    return Response.json({ 'status': 'error', 'message': 'Authorization code missing.' }, {
                        status: 400
                    });
                }

                const tokenRequest = await fetch('https://auth.tesla.com/oauth2/v3/token', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded'
                    },
                    body: new URLSearchParams({
                        'grant_type': 'authorization_code',
                        'client_id': env.CLIENT_ID,
                        'client_secret': env.CLIENT_SECRET,
                        'code': code,
                        'audience': env.AUDIENCE,
                        'redirect_uri': `https://${url.hostname}/callback`
                    })
                });
                const tokenResponse = await gatherResponse(tokenRequest);
                if (tokenRequest.status != 200) {
                    console.error(tokenResponse);
                    return Response.json({ 'status': 'error', 'message': 'Authorization code flow failed.' }, {
                        status: 500
                    });
                }

                const userRequest = await fetch(`${env.AUDIENCE}/api/1/users/me`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${tokenResponse.access_token}`
                    }
                });
                const userResponse = await gatherResponse(userRequest);
                if (userRequest.status != 200) {
                    console.error(userResponse);
                    return Response.json({ 'status': 'error', 'message': 'Failed to get user data.' }, {
                        status: 500
                    });
                }

                if (userResponse.response.email == env.ALLOWED_EMAIL) {
                    return Response.json(tokenResponse);
                } else {
                    return Response.json({ 'status': 'error', 'message': 'Tesla account is not authorized for this application.' }, {
                        status: 401
                    });
                }
            default:
                return Response.redirect(`https://auth.tesla.com/oauth2/v3/authorize?&client_id=${env.CLIENT_ID}&locale=en-US&prompt=login&redirect_uri=https%3A%2F%2F${url.hostname}%2Fcallback&response_type=code&scope=${encodeURI(env.SCOPE)}&state=mkb6_4Vvj5ERbqy-pUGSz`, 301);
        }
    }
};