// Error component.
Vue.component('api-secret', {
	props: ['secret'],
	template: '<div class="alert alert-info" v-if="secret">' +
		'<h4>Secret Info:</h4>' +
		'<ul v-for="data, key in secret">' +
		'<li><strong>{{ key }}</strong>: {{ data }}</li>' +
		'</ul>' +
		'</div>'
});