"use strict"
$(function() {
	$(".timeField").datetimepicker({
		format: "HH:mm",
		pickDate: false,
		language: "sk",

	});
	$("input[name=overrideState]").on("click", function(){
		if ($("input[name=overrideState]:checked").val() !== "0") {
			$("input:not(input[name=overrideState]):not(input[type=submit])").attr("disabled", "true");
		} else {
			$("input:not(input[name=overrideState]):not(input[type=submit])").removeAttr("disabled");
		}
	});
})