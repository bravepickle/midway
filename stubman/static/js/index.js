function initIndex() {
	$('#btn-create').click(function() {
		top.location.href = 'create'
	});
	
	var delUri
	$('#del-confirm').on('shown.bs.modal', function (ev) {
 		$('#btn-del-cancel').focus()
		delUri = $(ev.relatedTarget).attr('data-href')
	})
		
	$('#btn-del-confirm').click(function() {
		$.post('/stubman/delete/1', {}, function (resp) { 
			top.location.reload(); // reload page to see results
		});
	});
}