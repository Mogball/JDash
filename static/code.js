const urlData = $('#url-data');
$(function () {
    $('#submit-encode').on('click', function (e) {
        e.preventDefault();
        $.post(urlData.data('encode'), $('#message-form').serialize(), function (encoded) {
            $('#encoded').val(encoded);
        });
    });
    $('#submit-decode').on('click', function (e) {
        e.preventDefault();
        $.post(urlData.data('decode'), $('#encoded-form').serialize(), function (message) {
            $('#message').val(message);
        });
    });
});
