new Vue({
	el: '#set',
	data: {
		error_response: null, // API error response, if any.
		secret: null, // API secret response, if any.
		data: null, // Data to be put into the secret.
		expireAfterViews: null, // How many times secret can be viewed before it expires.
		expireAfter: null, // How many minutes does the secret exist for.
	},
	methods: {
		putSecret: function () {
			// We have to create request string manually since our API consumes URL-encoded requests while axios
			// sends data as json payload by default.
			const params = new URLSearchParams();
			params.append('secret', this.data);
			params.append('expireAfterViews', this.expireAfterViews);
			params.append('expireAfter', this.expireAfter);

			// Perform the request itself.
			axios
				.post('/secret/server/v1/secret/', params)
				.then(response => {
					this.secret = response.data;
					this.error_response = null;
				})
				.catch(response => {
					this.secret = null;
					this.error_response = response;
				})
		}
	}
});
