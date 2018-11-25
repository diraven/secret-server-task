// Error component.
Vue.component('api-error', {
	props: ['response'],
	template: '<div class="alert alert-danger" v-if="response">' +
		'<h4 v-if="response.response.statusText">{{response.response.status}} {{response.response.statusText}}</h4>' +
		'<div v-if="response.response.data.message">{{response.response.data.message}}</div>' +
		'</div>'
});