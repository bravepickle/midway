{{ define "title"}}<title>Stubman | Stub Edit #{{.Data.Id}}</title>{{ end }}
{{ define "content"}}<h1>Stub Edit #{{.Data.Id}}</h1>

<form method="POST" class="stub-form">
<fieldset>
<legend>Basic</legend>
  <div class="form-group">
    <label for="model-id">ID</label>
    <input readonly="readonly" type="text" class="form-control" value="{{.Data.Id}}" id="model-id" placeholder="ID">
  </div>
  <div class="form-group">
    <label for="model-name">Name</label>
    <input type="text" class="form-control" value="{{.Data.Name}}" id="model-name" name="name" placeholder="Name">
  </div>

  <div class="form-group">
    <label for="model-id">Created Date</label>
    <input readonly="readonly" type="text" class="form-control" value="{{.Data.Created}}" id="model-created" placeholder="Created Date">
  </div>
</fieldset>

<fieldset>
<legend>Request</legend>
  <div class="form-group">
    <label for="model-request_method">Request Method</label>
	<select  class="form-control" id="model-request_method" name="request_method">
	  <option {{ if eq .Data.RequestMethod `` }} selected="selected" {{end}}></option>
	  <option {{ if eq .Data.RequestMethod `GET` }} selected="selected" {{end}}>GET</option>
	  <option {{ if eq .Data.RequestMethod `POST` }} selected="selected" {{end}}>POST</option>
	  <option {{ if eq .Data.RequestMethod `PUT` }} selected="selected" {{end}}>PUT</option>
	  <option {{ if eq .Data.RequestMethod `PATCH` }} selected="selected" {{end}}>PATCH</option>
	  <option {{ if eq .Data.RequestMethod `DELETE` }} selected="selected" {{end}}>DELETE</option>
	  <option {{ if eq .Data.RequestMethod `HEAD` }} selected="selected" {{end}}>HEAD</option>
	  <option {{ if eq .Data.RequestMethod `OPTIONS` }} selected="selected" {{end}}>OPTIONS</option>
	</select>
  </div>
  <div class="form-group">
    <label for="model-request_uri">Request URI</label>
    <input type="text" class="form-control" value="{{.Data.RequestUri}}" id="model-request_uri" name="request_uri" placeholder="Request URI">
  </div>
  
<div class="headers-group form-group">
	<label for="request-headers">Request Headers</label>
	<div id="request-headers" class="headers-list">
	{{range .Data.RequestParsed.Headers}}
		<div>
			<input class="form-control" name="request[headers][]" placeholder="Content-Type: application/json" type=text value="{{.}}" /> 
			<a href="#1" class="btn-del-header"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></a>
		</div>
	{{end}}
	</div>
	
	<div class="add-header-container">
		<button type="button" id="request-add-header" class="btn btn-success btn-add-header">Add Header</button>
	</div>
</div>
<div class="form-group clearfix" id="request_body_container">
    <label for="model-request_body">Request Body</label>
	<div class="pull-right">
		<a class="decoder" href="#">expand JSON</a> | 
		<a class="encoder" href="#">fold JSON</a>
	</div>
    <textarea type="text" class="form-control" id="model-request_body" name="request[body]" placeholder="Request Body">{{.Data.RequestParsed.Body}}</textarea>
</div>
</fieldset>

<fieldset>
<legend>Response</legend>
<div class="headers-group form-group">
	<label for="response-headers">Response Headers</label>
	<div id="response-headers" class="headers-list">
	{{range .Data.ResponseParsed.Headers}}
		<div>
			<input class="form-control" name="request[headers][]" placeholder="Content-Type: application/json" type=text value="{{.}}" /> 
			<a href="#2" class="btn-del-header"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></a>
		</div>
	{{end}}
	</div>
	
	<div class="add-header-container">
		<button type="button" id="response-add-header" class="btn btn-success btn-add-header">Add Header</button>
	</div>
</div>
<div class="form-group clearfix">
    <label for="model-response_body">Response Body</label>
	<div class="pull-right">
		<a class="decoder" href="#">expand JSON</a> |  
		<a class="encoder" href="#">fold JSON</a>
	</div>
    <textarea type="text" class="form-control" id="model-response_body" name="response[body]" placeholder="Response Body">{{.Data.ResponseParsed.Body}}</textarea>
</div>
</fieldset>

<div class="clearfix">
	<button type="submit" class="btn btn-primary pull-right">Update</button>
</div>
</form>

<script type="text/template" id="header-request-tpl">
<div>
	<input class="form-control" name="request[headers][]" placeholder="Content-Type: application/json" type=text /> 
	<a href="#3" class="btn-del-header"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></a>
</div>
</script>

<script type="text/template" id="header-response-tpl">
<div>
	<input class="form-control" name="response[headers][]" placeholder="Content-Type: application/json" type=text /> 
	<a href="#4" class="btn-del-header"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></a>
</div>
</script>

{{ end }}

{{ define "scripts" }}
<script src="/stubman/static/js/edit.js" crossorigin="anonymous"></script>
<script type="text/javascript">
   $(document).ready(function() {
		initEdit();
   });
</script>
{{ end }}