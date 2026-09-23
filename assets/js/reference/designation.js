'use strict';

// Global DataTables configuration defaults
window.grid_limit = window.grid_limit || [[10, 20, 50, 100, 200], [10, 20, 50, 100, 200]];
window.grid_limit_default = window.grid_limit_default || 20;
window.grid_height = window.grid_height || 450;
window.actual_status = {
    1: "New Add",
    2: "Add-Edit",
    3: "Edit",
    4: "Delete",
    5: "Deactive",
    6: "Active"
};
window.process_status = {
    0: "",
    1: "(Sent to Checker)",
    2: "(Returned)",
    3: "(Approved)"
};
window.data_status = {
    0: "Deactive",
    1: "Active"
};
window.duty_carry_out = window.duty_carry_out || {};
if (typeof window.super_admin_status === 'undefined') {
    window.super_admin_status = 1;
}

function showLoading() {
    // Optional loading indicator hook
}

function response_error_alert() {
    if (typeof showAlertBox === 'function') {
        showAlertBox({ title: 'Error', message: 'Something went wrong fetching data', type: 'danger' });
    }
}

function grid_response(json) {
    if (json && Array.isArray(json.data)) return json.data;
    if (json && json.data && Array.isArray(json.data.data)) return json.data.data;
    if (json && Array.isArray(json.items)) return json.items;
    return [];
}

let selectedRowData = null;

// Global get_new_page dispatcher
window.get_new_page = function (url, title, id, name, popupStatus, width) {
    id = id ? String(id).trim() : '';
    name = name ? decodeURIComponent(name) : '';

    if (!id) {
        // Triggered from top action buttons
        if (url.indexOf('/add') !== -1) {
            openAddModal();
            return;
        }
        if (!selectedRowData) {
            if (typeof showAlertBox === 'function') {
                showAlertBox({ title: 'Notice', message: 'Please select a row from the table first', type: 'warning' });
            } else {
                alert('Please select a row from the table first');
            }
            return;
        }
        id = selectedRowData.uuid || selectedRowData.id;
        name = selectedRowData.designation_name || selectedRowData.name;
    }

    if (url.indexOf('/view') !== -1) {
        openViewModal(id);
    } else if (url.indexOf('/edit') !== -1) {
        openEditModal(id);
    } else if (url.indexOf('/delete') !== -1) {
        openDeleteModal(id, name);
    } else if (url.indexOf('/stc') !== -1) {
        openStcModal(id, name);
    } else if (url.indexOf('/verify') !== -1) {
        openVerifyModal(id, name);
    } else if (url.indexOf('/active') !== -1) {
        openStatusModal('active', id, name);
    } else if (url.indexOf('/deactivate') !== -1 || url.indexOf('/deactive') !== -1) {
        openStatusModal('deactivate', id, name);
    }
};

window.return_message_history = function (name, id, table) {
    name = name ? decodeURIComponent(name) : 'Designation';
    openHistoryModal(id, name);
};

function openAddModal() {
    $('#form_designation')[0].reset();
    $('#designation_uuid').val('');
    $('#designation_id').val('');
    $('#modal_designation_title').text('Add Designation');
    $('#btn_save_designation').text('Save Designation');
    $('#modal_designation_form').modal('show');
}

function openEditModal(targetId) {
    $.get(`/reference/designation/edit/${targetId}`)
        .done(function (res) {
            let data = (res && res.data) || {};
            $('#form_designation')[0].reset();
            $('#designation_uuid').val(data.uuid || targetId);
            $('#designation_id').val(data.id || '');
            $('#designation_name').val(data.name || data.designation_name || '');
            $('#designation_code').val(data.code || data.designation_code || '');
            $('#designation_title').val(data.title || data.designation_title || '');
            $('#designation_details').val(data.details || data.designation_details || '');
            $('#modal_designation_title').text('Edit Designation');
            $('#btn_save_designation').text('Update Designation');
            $('#modal_designation_form').modal('show');
        })
        .fail(function (xhr) {
            let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Failed to fetch designation data';
            if (typeof showAlertBox === 'function') showAlertBox({ title: 'Error', message: msg, type: 'danger' });
            else alert(msg);
        });
}

function openViewModal(targetId) {
    $.get(`/reference/designation/view/${targetId}`)
        .done(function (res) {
            let d = (res && res.data) || {};
            $('#view_name').text(d.name || d.designation_name || '-');
            $('#view_code').text(d.code || d.designation_code || '-');
            $('#view_title').text(d.title || d.designation_title || '-');
            $('#view_details').text(d.details || d.designation_details || '-');

            let opText = window.actual_status[d.actual_status] || 'New Add';
            $('#view_op_status').text(opText);
            $('#view_status').text(d.status === 1 ? 'Active' : 'Deactive');

            let procText = 'Draft';
            if (d.process_status === 1) procText = 'Sent to Checker';
            else if (d.process_status === 2) procText = 'Returned';
            else if (d.process_status === 3) procText = 'Approved';
            $('#view_process_status').text(procText);

            $('#view_checker').text(d.checker > 0 ? `User ID ${d.checker}` : 'None');
            $('#view_message').text(d.message || '-');
            $('#view_remarks').text(d.remarks || '-');
            $('#view_created_at').text(d.created_at || '-');

            $('#btn_open_history_from_view').data('uuid', d.uuid || targetId);
            $('#modal_designation_view').modal('show');
        })
        .fail(function (xhr) {
            let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Failed to fetch details';
            if (typeof showAlertBox === 'function') showAlertBox({ title: 'Error', message: msg, type: 'danger' });
            else alert(msg);
        });
}

function openDeleteModal(targetId, name) {
    $('#form_designation_delete')[0].reset();
    $('#delete_target_id').val(targetId);
    $('#delete_designation_name').text(name || 'Selected Item');
    $('#modal_designation_delete').modal('show');
}

function openStcModal(targetId, name) {
    $('#form_designation_stc')[0].reset();
    $('#stc_target_id').val(targetId);
    $('#stc_designation_name').text(name || 'Selected Item');
    $('#modal_designation_stc').modal('show');
}

function openVerifyModal(targetId, name) {
    $('#form_designation_verify')[0].reset();
    $('#verify_target_id').val(targetId);
    $('#verify_designation_name').text(name || 'Selected Item');
    $('#verify_type_approve').prop('checked', true);
    $('#group_return_reason').addClass('d-none');
    $('#modal_designation_verify').modal('show');
}

function openStatusModal(actionType, targetId, name) {
    $('#form_designation_status')[0].reset();
    $('#status_target_id').val(targetId);
    $('#status_action_type').val(actionType);
    $('#status_designation_name').text(name || 'Selected Item');

    let isActivate = actionType === 'active';
    $('#modal_status_title').text(isActivate ? 'Activate Designation' : 'Deactivate Designation');
    $('#status_prompt_text').text(isActivate ? 'Reason for activating this designation:' : 'Reason for deactivating this designation:');
    $('#btn_confirm_status').text(isActivate ? 'Activate' : 'Deactivate');
    $('#modal_designation_status').modal('show');
}

function openHistoryModal(targetId, name) {
    $.get(`/reference/designation/history/${targetId}`)
        .done(function (res) {
            let list = (res && res.data) || [];
            let html = '';
            if (!list || list.length === 0) {
                html = '<tr><td colspan="9" class="text-center py-3 text-muted">No audit history recorded yet</td></tr>';
            } else {
                list.forEach(function (h, idx) {
                    let op = window.actual_status[h.actual_status] || 'New Add';
                    let ps = h.process_status === 3 ? 'Approved' : (h.process_status === 1 ? 'Sent to Checker' : (h.process_status === 2 ? 'Returned' : 'Draft'));
                    let st = h.status === 1 ? 'Active' : 'Deactive';

                    html += `
                        <tr>
                            <td>${idx + 1}</td>
                            <td>${h.name || h.designation_name || '-'}</td>
                            <td>${h.code || h.designation_code || '-'}</td>
                            <td>${h.title || h.designation_title || '-'}</td>
                            <td><span class="badge bg-light text-dark border">${op}</span></td>
                            <td><span class="badge bg-primary-transparent text-primary">${ps}</span></td>
                            <td><span class="badge ${h.status === 1 ? 'bg-success-transparent text-success' : 'bg-danger-transparent text-danger'}">${st}</span></td>
                            <td>${h.message || '-'}</td>
                            <td>${h.created_at || '-'}</td>
                        </tr>
                    `;
                });
            }
            $('#history_table_body').html(html);
            $('#modal_designation_history').modal('show');
        })
        .fail(function (xhr) {
            let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Failed to load history';
            if (typeof showAlertBox === 'function') showAlertBox({ title: 'Error', message: msg, type: 'danger' });
            else alert(msg);
        });
}

function reloadTable() {
    if ($.fn.DataTable.isDataTable('#designation_data_table')) {
        $('#designation_data_table').DataTable().ajax.reload(null, false);
    }
}

$(document).ready(function () {
    let gridRight = window.grid_right || [];
    let columns = [];
    let nonOrderableTargets = [];

    // Build right action columns matching the user's Blade logic
    gridRight.forEach(function (right, idx) {
        nonOrderableTargets.push(idx);
        columns.push({
            data: "id",
            render: function (data, type, full, meta) {
                var titleValue = encodeURIComponent(full.designation_name || full.name || '');
                var targetId = full.uuid || full.id;
                var action = right.r_action_name;
                var iconHtml = '<i class="' + right.r_icon + ' fs-16"></i>';

                if (action === 'view') {
                    return '<a href="#" class="text-info me-1" title="' + right.r_title + '" onclick="get_new_page(\'' + right.r_route_name + '\',\'' + right.r_title + '\',\'' + targetId + '\',\'' + titleValue + '\',\'' + right.popup_status + '\',\'' + right.width + '\');" >' + iconHtml + '</a>';
                }

                if (action === 'edit') {
                    if ((full.process_status == 2 || full.process_status == 3 || full.process_status == '0' || full.process_status == '')) {
                        if (full.process_status == 3 || full.process_status == 2) {
                            return '<a href="#" class="text-primary me-1" title="' + right.r_title + '" onclick="get_new_page(\'' + right.r_route_name + '\',\'' + right.r_title + '\',\'' + targetId + '\',\'' + titleValue + '\',\'' + right.popup_status + '\',\'' + right.width + '\');" >' + iconHtml + '</a>';
                        } else if (full.actual_status != 4) {
                            return '<a href="#" class="text-primary me-1" title="' + right.r_title + '" onclick="get_new_page(\'' + right.r_route_name + '\',\'' + right.r_title + '\',\'' + targetId + '\',\'' + titleValue + '\',\'' + right.popup_status + '\',\'' + right.width + '\');" >' + iconHtml + '</a>';
                        }
                    }
                    return '';
                }

                if (action === 'delete') {
                    if ((full.process_status == 2 || full.process_status == 3 || full.process_status == '0' || full.process_status == '')) {
                        if (full.process_status == 3 || full.process_status == 2) {
                            return '<a href="#" class="text-danger me-1" title="' + right.r_title + '" onclick="get_new_page(\'' + right.r_route_name + '\',\'' + right.r_title + '\',\'' + targetId + '\',\'' + titleValue + '\',\'' + right.popup_status + '\',\'' + right.width + '\');" >' + iconHtml + '</a>';
                        } else if (full.actual_status != 3) {
                            return '<a href="#" class="text-danger me-1" title="' + right.r_title + '" onclick="get_new_page(\'' + right.r_route_name + '\',\'' + right.r_title + '\',\'' + targetId + '\',\'' + titleValue + '\',\'' + right.popup_status + '\',\'' + right.width + '\');" >' + iconHtml + '</a>';
                        }
                    }
                    return '';
                }

                if (action === 'sent_to_checker') {
                    if (full.process_status != 3 && full.process_status != 1 && full.actual_status != 5 && full.actual_status != 6) {
                        return '<a href="#" class="text-primary me-1" title="' + right.r_title + '" onclick="get_new_page(\'' + right.r_route_name + '\',\'' + right.r_title + '\',\'' + targetId + '\',\'' + titleValue + '\',\'' + right.popup_status + '\',\'' + right.width + '\');" >' + iconHtml + '</a>';
                    }
                    return '';
                }

                if (action === 'verify') {
                    var canVerify = (typeof window.duty_carry_out[full.checker] !== 'undefined' || 1 == window.super_admin_status);
                    if (full.process_status == 1 && canVerify) {
                        return '<a href="#" class="text-success me-1" title="' + right.r_title + '" onclick="get_new_page(\'' + right.r_route_name + '\',\'' + right.r_title + '\',\'' + targetId + '\',\'' + titleValue + '\',\'' + right.popup_status + '\',\'' + right.width + '\');" >' + iconHtml + '</a>';
                    }
                    return '';
                }

                if (action === 'active') {
                    if ((full.process_status == 3 || full.process_status == 2 || full.actual_status == 5) && full.status == 0) {
                        return '<a href="#" class="text-success me-1" title="' + right.r_title + '" onclick="get_new_page(\'' + right.r_route_name + '\',\'' + right.r_title + '\',\'' + targetId + '\',\'' + titleValue + '\',\'' + right.popup_status + '\',\'' + right.width + '\');" >' + iconHtml + '</a>';
                    }
                    return '';
                }

                if (action === 'deactive') {
                    if ((full.process_status == 3 || full.process_status == 2 || full.actual_status == 6) && full.status == 1) {
                        return '<a href="#" class="text-warning me-1" title="' + right.r_title + '" onclick="get_new_page(\'' + right.r_route_name + '\',\'' + right.r_title + '\',\'' + targetId + '\',\'' + titleValue + '\',\'' + right.popup_status + '\',\'' + right.width + '\');" >' + iconHtml + '</a>';
                    }
                    return '';
                }

                return '';
            }
        });
    });

    // Data columns matching user's specification
    columns.push(
        {
            data: "designation_name",
            render: function (data, type, full, meta) {
                return '<div class="text-start fw-semibold text-dark">' + (data || full.name || '-') + '</div>';
            }
        },
        {
            data: "designation_code",
            render: function (data, type, full, meta) {
                var codeVal = data || full.code;
                return '<div class="text-start">' + (codeVal ? '<span class="badge bg-light text-dark border">' + codeVal + '</span>' : '-') + '</div>';
            }
        },
        {
            data: "actual_status",
            render: function (data, type, full, meta) {
                var statusText = '';
                if (typeof window.process_status[full.process_status] !== 'undefined') {
                    statusText = ' ' + window.process_status[full.process_status];
                }
                var actLabel = window.actual_status[data] || (window.actual_status[full.actual_status] || 'New Add');
                return '<div class="text-start"><span class="badge bg-primary-transparent text-primary fw-medium">' + actLabel + statusText + '</span></div>';
            }
        },
        {
            data: "status",
            render: function (data, type, full, meta) {
                var st = data !== undefined ? data : full.status;
                var badgeClass = st == 1 ? 'bg-success-transparent text-success' : 'bg-danger-transparent text-danger';
                return '<div class="text-start"><span class="badge ' + badgeClass + ' fw-semibold">' + (window.data_status[st] || (st == 1 ? 'Active' : 'Deactive')) + '</span></div>';
            }
        },
        {
            data: "message",
            render: function (data, type, full, meta) {
                if (data != null && data !== '') {
                    var safeName = encodeURIComponent(full.designation_name || full.name || '');
                    var safeId = full.uuid || full.id;
                    return '<div class="text-start cursor-pointer text-primary" onclick="return_message_history(\'' + safeName + '\',\'' + safeId + '\',\'designations\')">' + data + '</div>';
                }
                return '';
            }
        }
    );

    let defaultOrderCol = gridRight.length; // Order on designation_name

    let table = $('#designation_data_table').DataTable({
        lengthMenu: window.grid_limit,
        pageLength: window.grid_limit_default,
        serverSide: true,
        responsive: true,
        colReorder: true,
        scrollX: true,
        autoWidth: true,
        scrollY: window.grid_height + "px",
        scrollCollapse: false,
        ajax: {
            url: "/reference/designation/view",
            type: "POST",
            beforeSend: function (xhr) {
                var csrf = $('meta[name="csrf-token"]').attr('content') || $('input[name="_csrf"]').val();
                if (csrf) xhr.setRequestHeader('X-CSRF-Token', csrf);
                showLoading();
            },
            error: function (xhr, status, error) {
                response_error_alert();
            },
            dataSrc: function (json) {
                return grid_response(json);
            }
        },
        order: [
            [defaultOrderCol, 'desc']
        ],
        columns: columns,
        columnDefs: [
            {
                targets: nonOrderableTargets,
                searchable: false,
                orderable: false
            }
        ]
    });

    // Row selection listener
    $('#designation_data_table tbody').on('click', 'tr', function () {
        if ($(this).hasClass('selected')) {
            $(this).removeClass('selected table-active');
            selectedRowData = null;
        } else {
            table.$('tr.selected').removeClass('selected table-active');
            $(this).addClass('selected table-active');
            selectedRowData = table.row(this).data();
        }
    });

    // Verify reason display toggle
    $('input[name="verify_type"]').on('change', function () {
        if ($(this).val() === '0') {
            $('#group_return_reason').removeClass('d-none');
            $('#verify_return_reason').prop('required', true);
        } else {
            $('#group_return_reason').addClass('d-none');
            $('#verify_return_reason').prop('required', false);
        }
    });

    $('#btn_open_history_from_view').on('click', function () {
        let uuid = $(this).data('uuid');
        $('#modal_designation_view').modal('hide');
        openHistoryModal(uuid, 'Designation');
    });

    // Add / Edit form submit
    $('#form_designation').on('submit', function (e) {
        e.preventDefault();
        let uuid = $('#designation_uuid').val();
        let url = uuid ? `/reference/designation/edit/${uuid}` : '/reference/designation/add';
        let formData = $(this).serialize();

        $.post(url, formData)
            .done(function (res) {
                $('#modal_designation_form').modal('hide');
                if (typeof showAlertBox === 'function') {
                    showAlertBox({ title: 'Success', message: res.message || 'Saved successfully', type: 'theme' });
                }
                reloadTable();
            })
            .fail(function (xhr) {
                let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Operation failed';
                if (typeof showAlertBox === 'function') showAlertBox({ title: 'Warning', message: msg, type: 'warning' });
                else alert(msg);
            });
    });

    // Delete form submit
    $('#form_designation_delete').on('submit', function (e) {
        e.preventDefault();
        let uuid = $('#delete_target_id').val();
        let formData = $(this).serialize();

        $.post(`/reference/designation/delete/${uuid}`, formData)
            .done(function (res) {
                $('#modal_designation_delete').modal('hide');
                if (typeof showAlertBox === 'function') {
                    showAlertBox({ title: 'Success', message: res.message || 'Deleted successfully', type: 'theme' });
                }
                reloadTable();
            })
            .fail(function (xhr) {
                let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Delete failed';
                if (typeof showAlertBox === 'function') showAlertBox({ title: 'Error', message: msg, type: 'danger' });
                else alert(msg);
            });
    });

    // Sent to checker submit
    $('#form_designation_stc').on('submit', function (e) {
        e.preventDefault();
        let uuid = $('#stc_target_id').val();
        let formData = $(this).serialize();

        $.post(`/reference/designation/stc/${uuid}`, formData)
            .done(function (res) {
                $('#modal_designation_stc').modal('hide');
                if (typeof showAlertBox === 'function') {
                    showAlertBox({ title: 'Success', message: res.message || 'Sent To Checker successfully', type: 'theme' });
                }
                reloadTable();
            })
            .fail(function (xhr) {
                let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Action failed';
                if (typeof showAlertBox === 'function') showAlertBox({ title: 'Error', message: msg, type: 'danger' });
                else alert(msg);
            });
    });

    // Verify submit
    $('#form_designation_verify').on('submit', function (e) {
        e.preventDefault();
        let uuid = $('#verify_target_id').val();
        let formData = $(this).serialize();

        $.post(`/reference/designation/verify/${uuid}`, formData)
            .done(function (res) {
                $('#modal_designation_verify').modal('hide');
                if (typeof showAlertBox === 'function') {
                    showAlertBox({ title: 'Success', message: res.message || 'Verified successfully', type: 'theme' });
                }
                reloadTable();
            })
            .fail(function (xhr) {
                let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Verification failed';
                if (typeof showAlertBox === 'function') showAlertBox({ title: 'Error', message: msg, type: 'danger' });
                else alert(msg);
            });
    });

    // Active / Deactivate form submit
    $('#form_designation_status').on('submit', function (e) {
        e.preventDefault();
        let uuid = $('#status_target_id').val();
        let actionType = $('#status_action_type').val();
        let formData = $(this).serialize();

        let url = actionType === 'active'
            ? `/reference/designation/active/${uuid}`
            : `/reference/designation/deactivate/${uuid}`;

        $.post(url, formData)
            .done(function (res) {
                $('#modal_designation_status').modal('hide');
                if (typeof showAlertBox === 'function') {
                    showAlertBox({ title: 'Success', message: res.message || 'Status updated successfully', type: 'theme' });
                }
                reloadTable();
            })
            .fail(function (xhr) {
                let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Status update failed';
                if (typeof showAlertBox === 'function') showAlertBox({ title: 'Error', message: msg, type: 'danger' });
                else alert(msg);
            });
    });
});
