'use strict';

(function ($) {
    let designationGrid = null;
    let selectedDesignation = null;

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

    $(document).ready(function () {
        initDesignationGrid();
        bindFormSubmissions();
    });

    function showNotify(type, message) {
        if (typeof showAlertBox === 'function') {
            let title = type === 'success' ? 'Success' : (type === 'warning' ? 'Warning' : 'Error');
            showAlertBox({
                title: title,
                message: message,
                type: type === 'danger' ? 'danger' : (type === 'warning' ? 'warning' : 'theme')
            });
        } else if (typeof toastNotification === 'function') {
            toastNotification(type, message);
        } else {
            alert(message);
        }
    }

    function initDesignationGrid() {
        designationGrid = initAgDataGrid({
            gridContainer: '#designation_ag_grid',
            url: '/reference/designation/view',
            exportFileName: 'designations_list',
            defaultOrderCol: 'id',
            defaultOrderDir: 'desc',
            totalBadgeSelector: '#designation_total_badge',
            rowSelection: 'single',
            onSelectionChanged: function (params) {
                if (params && params.api) {
                    let selectedRows = params.api.getSelectedRows();
                    selectedDesignation = selectedRows && selectedRows.length > 0 ? selectedRows[0] : null;
                }
            },
            columnDefs: [
                {
                    headerName: '#',
                    isIndexCol: true,
                    width: 70,
                    minWidth: 65
                },
                {
                    headerName: 'Actions',
                    field: 'actions',
                    pinned: 'left',
                    width: 175,
                    minWidth: 165,
                    sortable: false,
                    filter: false,
                    cellRenderer: function (params) {
                        let full = params.data;
                        if (!full) return '';
                        let titleValue = encodeURIComponent(full.designation_name || full.name || '');
                        let targetId = full.uuid || full.id;
                        let out = '<div class="d-flex align-items-center">';

                        // View
                        out += '<a href="#" class="text-info me-2" title="View" onclick="get_new_page(\'/reference/designation/view\',\'View Designation\',\'' + targetId + '\',\'' + titleValue + '\',1,1000);"><i class="ti ti-eye fs-16"></i></a>';

                        // Edit
                        if ((full.process_status == 2 || full.process_status == 3 || full.process_status == '0' || full.process_status == '')) {
                            if (full.process_status == 3 || full.process_status == 2 || full.actual_status != 4) {
                                out += '<a href="#" class="text-primary me-2" title="Edit" onclick="get_new_page(\'/reference/designation/edit\',\'Edit Designation\',\'' + targetId + '\',\'' + titleValue + '\',1,1000);"><i class="ti ti-edit fs-16"></i></a>';
                            }
                        }

                        // Delete
                        if ((full.process_status == 2 || full.process_status == 3 || full.process_status == '0' || full.process_status == '')) {
                            if (full.process_status == 3 || full.process_status == 2 || full.actual_status != 3) {
                                out += '<a href="#" class="text-danger me-2" title="Delete" onclick="get_new_page(\'/reference/designation/delete\',\'Delete Designation\',\'' + targetId + '\',\'' + titleValue + '\',1,1000);"><i class="ti ti-trash fs-16"></i></a>';
                            }
                        }

                        // Sent To Checker
                        if (full.process_status != 3 && full.process_status != 1 && full.actual_status != 5 && full.actual_status != 6) {
                            out += '<a href="#" class="text-primary me-2" title="Sent To Checker" onclick="get_new_page(\'/reference/designation/stc\',\'Sent To Checker\',\'' + targetId + '\',\'' + titleValue + '\',1,1000);"><i class="ti ti-arrow-forward-up fs-16"></i></a>';
                        }

                        // Verify
                        let canVerify = (typeof window.duty_carry_out[full.checker] !== 'undefined' || 1 == window.super_admin_status);
                        if (full.process_status == 1 && canVerify) {
                            out += '<a href="#" class="text-success me-2" title="Verify" onclick="get_new_page(\'/reference/designation/verify\',\'Verify Designation\',\'' + targetId + '\',\'' + titleValue + '\',1,1000);"><i class="ti ti-check fs-16"></i></a>';
                        }

                        // Active
                        if ((full.process_status == 3 || full.process_status == 2 || full.actual_status == 5) && full.status == 0) {
                            out += '<a href="#" class="text-success me-2" title="Activate" onclick="get_new_page(\'/reference/designation/active\',\'Active Designation\',\'' + targetId + '\',\'' + titleValue + '\',1,1000);"><i class="ti ti-lock-open fs-16"></i></a>';
                        }

                        // Deactive
                        if ((full.process_status == 3 || full.process_status == 2 || full.actual_status == 6) && full.status == 1) {
                            out += '<a href="#" class="text-warning me-2" title="Deactivate" onclick="get_new_page(\'/reference/designation/deactivate\',\'Deactivate Designation\',\'' + targetId + '\',\'' + titleValue + '\',1,1000);"><i class="ti ti-lock fs-16"></i></a>';
                        }

                        out += '</div>';
                        return out;
                    }
                },
                {
                    field: 'name',
                    colId: 'name',
                    headerName: 'Designation Name',
                    width: 250,
                    minWidth: 200,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return `<span class="fw-semibold text-dark"><i class="ti ti-briefcase me-1 text-primary"></i>${params.value || (params.data && params.data.designation_name) || '-'}</span>`;
                    }
                },
                {
                    field: 'code',
                    colId: 'code',
                    headerName: 'Code',
                    width: 120,
                    minWidth: 100,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        let val = params.value || (params.data && params.data.designation_code);
                        return val ? `<span class="badge bg-light text-dark border">${val}</span>` : '<span class="text-muted">-</span>';
                    }
                },
                {
                    field: 'title',
                    colId: 'title',
                    headerName: 'Title',
                    width: 200,
                    minWidth: 160,
                    sortable: true,
                    filter: true
                },
                {
                    field: 'actual_status',
                    colId: 'actual_status',
                    headerName: 'OP-Status',
                    width: 180,
                    minWidth: 150,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        let full = params.data || {};
                        let statusText = '';
                        if (typeof window.process_status[full.process_status] !== 'undefined') {
                            statusText = ' ' + window.process_status[full.process_status];
                        }
                        let actLabel = window.actual_status[params.value] || 'New Add';
                        let badgeClass = 'bg-primary-transparent text-primary';
                        if (params.value === 4) badgeClass = 'bg-danger-transparent text-danger';
                        else if (params.value === 5) badgeClass = 'bg-dark-transparent text-dark';
                        else if (params.value === 6) badgeClass = 'bg-success-transparent text-success';
                        else if (params.value === 2) badgeClass = 'bg-warning-transparent text-warning';
                        return `<span class="badge ${badgeClass} fw-semibold">${actLabel}${statusText}</span>`;
                    }
                },
                {
                    field: 'status',
                    colId: 'status',
                    headerName: 'Status',
                    width: 120,
                    minWidth: 100,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        let st = params.value;
                        let badgeClass = st === 1 ? 'bg-success-transparent text-success' : 'bg-danger-transparent text-danger';
                        return `<span class="badge ${badgeClass} fw-semibold">${window.data_status[st] || (st === 1 ? 'Active' : 'Deactive')}</span>`;
                    }
                },
                {
                    field: 'message',
                    colId: 'message',
                    headerName: 'Message',
                    width: 220,
                    minWidth: 180,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        let full = params.data || {};
                        let msg = params.value;
                        if (!msg) return '<span class="text-muted">-</span>';
                        let safeName = encodeURIComponent(full.designation_name || full.name || '');
                        let safeId = full.uuid || full.id;
                        return `<span class="text-primary cursor-pointer text-truncate d-inline-block mw-100" onclick="return_message_history('${safeName}','${safeId}','designations')" title="${msg}">${msg}</span>`;
                    }
                },
                {
                    field: 'created_at',
                    colId: 'created_at',
                    headerName: 'Created At',
                    width: 180,
                    minWidth: 160,
                    sortable: true,
                    filter: true,
                    filterType: 'date',
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12"><i class="ti ti-clock me-1"></i>${params.value || '-'}</span>`;
                    }
                }
            ]
        });
    }

    function refreshGrid() {
        if (designationGrid && designationGrid.gridApi) {
            designationGrid.gridApi.refreshServerSide({ purge: true });
        }
    }

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
            if (!selectedDesignation) {
                showNotify('warning', 'Please select a row from the grid first');
                return;
            }
            id = selectedDesignation.uuid || selectedDesignation.id;
            name = selectedDesignation.designation_name || selectedDesignation.name;
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
                showNotify('danger', msg);
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
                showNotify('danger', msg);
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
                showNotify('danger', msg);
            });
    }

    function bindFormSubmissions() {
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
                    showNotify('success', res.message || 'Saved successfully');
                    refreshGrid();
                })
                .fail(function (xhr) {
                    let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Operation failed';
                    showNotify('warning', msg);
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
                    showNotify('success', res.message || 'Deleted successfully');
                    refreshGrid();
                })
                .fail(function (xhr) {
                    let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Delete failed';
                    showNotify('danger', msg);
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
                    showNotify('success', res.message || 'Sent To Checker successfully');
                    refreshGrid();
                })
                .fail(function (xhr) {
                    let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Action failed';
                    showNotify('danger', msg);
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
                    showNotify('success', res.message || 'Verified successfully');
                    refreshGrid();
                })
                .fail(function (xhr) {
                    let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Verification failed';
                    showNotify('danger', msg);
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
                    showNotify('success', res.message || 'Status updated successfully');
                    refreshGrid();
                })
                .fail(function (xhr) {
                    let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Status update failed';
                    showNotify('danger', msg);
                });
        });
    }
})(jQuery);
