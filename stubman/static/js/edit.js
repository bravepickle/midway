function initEdit() {
	var reqHeader = $('#header-request-tpl').html()
	var reqList = $('#request-headers')
	$('#request-add-header').click(function(e) {
		e.preventDefault();
		reqList.append(reqHeader)
	});
	
	var respHeader = $('#header-response-tpl').html()
	var respList = $('#response-headers')
	$('#response-add-header').click(function(e) {
		e.preventDefault();
		respList.append(respHeader)
	});
	
	$('.decoder').click(function(e) {
		e.preventDefault();
		var el = $(this).closest('.form-group').find('textarea');
		try {
			var parsed = JSON.parse(el.val().trim())
			el.val(JSON.stringify(parsed, null, '  '))	
		} catch(e) {
			console.log(e)
		}
	});
	
	$('.encoder').click(function(e) {
		e.preventDefault();
		var el = $(this).closest('.form-group').find('textarea');
		try {
			var parsed = JSON.parse(el.val().trim())
			el.val(JSON.stringify(parsed))	
		} catch(e) {
			console.log(e)
		}
	});
	
	$('.headers-group').on('click', '.btn-del-header', function(e) {
		e.preventDefault();
		$(this).parent().remove()
	});	
}
