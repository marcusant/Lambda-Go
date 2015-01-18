$(document).ready(function() {
    $('#themeselect').change(function() {
        var theme_name = $(this).find('option:selected').val();
        window.location.href = '/usercp?theme='+theme_name;
    })
});