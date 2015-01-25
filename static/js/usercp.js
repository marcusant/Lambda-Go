$(document).ready(function() {
    $('#themeselect').change(function() {
        var theme_name = $(this).find('option:selected').val();
        window.location.href = '/settheme?name='+theme_name;
    })
});