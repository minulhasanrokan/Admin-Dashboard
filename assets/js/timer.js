(function () {
    "use strict";

    $('.digit-group').find('input').each(function () {
        $(this).on('keyup', function (e) {
            let parent = $($(this).parent());
            if (e.keyCode === 8 || e.keyCode === 37) {
                let prev = parent.find('input#' + $(this).data('previous'));
                if (prev.length) {
                    $(prev).select();
                }
            } else if ((e.keyCode >= 48 && e.keyCode <= 57) || (e.keyCode >= 65 && e.keyCode <= 90) || (e.keyCode >= 96 && e.keyCode <= 105) || e.keyCode === 39) {
                let next = parent.find('input#' + $(this).data('next'));
                if (next.length) {
                    $(next).select();
                } else {
                    if (parent.data('autosubmit')) {
                        parent.submit();
                    }
                }
            }
        });
    });

    $('.digit-group input').on('keyup', function () {
        let self = $(this);
        if (self.val() !== '') {
            self.addClass('active');
        } else {
            self.removeClass('active');
        }
    });
})();

function startTimer() {
    updateTimer();
    return setInterval(updateTimer, 1000);
}

function updateTimer() {
    let min = Math.floor(time / 60);
    let sec = time % 60;

    timerEl.textContent =
        String(min).padStart(2, '0') + ":" + String(sec).padStart(2, '0');

    if (time > 0) {
        time--;
    } else {
        clearInterval(interval);

        timerEl.textContent = "00:00";

        resendBtn.classList.remove("disabled-link");
    }
}

document.getElementById("digit_one").addEventListener("focusout", function () {

    form_validation({control: "digit_one#number#mandatory", msg_text: "Digit One"});
});

document.getElementById("digit_two").addEventListener("focusout", function () {

    form_validation({control: "digit_two#number#mandatory", msg_text: "Digit Two"});
});

document.getElementById("digit_three").addEventListener("focusout", function () {

    form_validation({control: "digit_three#number#mandatory", msg_text: "Digit Three"});
});

document.getElementById("digit_four").addEventListener("focusout", function () {

    form_validation({control: "digit_four#number#mandatory", msg_text: "Digit Four"});
});