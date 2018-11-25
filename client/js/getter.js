new Vue({
	el: '#get',
	data: {
		error_response: null, // API error response if any.
		secret: null, // API secret response if any.
		hash: null, // Hash for the secret to be retrieved by.
	},
	methods: {
		getSecret: function () {
			// Perform the get request.
			axios
				.get('/secret/server/v1/secret/' + this.hash)
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