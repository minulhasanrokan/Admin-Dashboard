function confirmDialog(title, text, icon = "question") {
    return Swal.fire({
        title: title,
        text: text,
        icon: icon,
        showCancelButton: true,
        confirmButtonText: "OK",
        cancelButtonText: "Cancel",
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        reverseButtons: true,
        allowOutsideClick: false,
        allowEscapeKey: false
    }).then((result) => {
        return result.isConfirmed;
    });
}

function showSuccess(title, text) {
    return Swal.fire({
        title: title,
        text: text,
        icon: "success",
        confirmButtonText: "OK"
    });
}

function showError(title, text) {
    return Swal.fire({
        title: title,
        text: text,
        icon: "error",
        confirmButtonText: "OK"
    });
}

function showWarning(title, text) {
    return Swal.fire({
        title: title,
        text: text,
        icon: "warning",
        confirmButtonText: "OK"
    });
}

function showInfo(title, text) {
    return Swal.fire({
        title: title,
        text: text,
        icon: "info",
        confirmButtonText: "OK"
    });
}