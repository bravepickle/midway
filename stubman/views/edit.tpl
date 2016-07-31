{{ define "title"}}<title>Stubman | Stub Edit #{{.Id}}</title>{{ end }}
{{ define "content"}}<h1>Stub Edit #{{.Id}}</h1>

<div>ID: <input type=text value="{{.Id}}" readonly /></div>
<div>Name: <input type=text value="{{.Name}}" /></div>
<div>Method: <input type=text value="{{.RequestMethod}}" /></div>
<div>URI: <input type=text value="{{.RequestUri}}" /></div>
<div>Request Headers:</div>

{{range .RequestParsed.Headers}}
<div><input type=text value="{{.}}" /> <button class="btn btn-danger">Del</button></div>
{{end}}


<script type="text/javascript">
   $(document).ready(function() {
		$('#btn-create').click(function() {
			top.location.href = 'create'
		});
		$('.btn-del').click(function(el) {
			console.log(this)
			console.log(el)
		});
		
		$('#del-confirm').on('shown.bs.modal', function () {
  			$('#btn-del-cancel').focus()
		})
 		
   });
</script>
{{ end }}
