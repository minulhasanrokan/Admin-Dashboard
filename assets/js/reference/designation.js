'use strict';

(function ($) {
    let designationGrid = null;
    let selectedDesignation = null;

    $(document).ready(function () {
        initDesignationGrid();
        bindTopActions();
        bindRowActions();
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
                    headerName: 'Actions',
                    field: 'actions',
                    pinned: 'left',
                    width: 175,
                    minWidth: 165,
                    sortable: false,
                    filter: false,
                    cellRenderer: function (params) {
                        let row = params.data;
                        if (!row) return '';
                        let targetId = row.uuid || row.id;
                        let safeName = encodeURIComponent(row.name || '');

                        let statusIcon = row.status === 1
                            ? `<button type="button" class="btn btn-link text-warning p-0 me-2 btn-action-deactive" title="Deactivate" data-uuid="${targetId}" data-name="${safeName}"><i class="ti ti-lock fs-16"></i></button>`
                            : `<button type="button" class="btn btn-link text-success p-0 me-2 btn-action-active" title="Activate" data-uuid="${targetId}" data-name="${safeName}"><i class="ti ti-lock-open fs-16"></i></button>`;

                        return `
                            <div class="d-flex align-items-center">
                                <button type="button" class="btn btn-link text-danger p-0 me-2 btn-action-delete" title="Delete" data-uuid="${targetId}" data-name="${safeName}"><i class="ti ti-trash fs-16"></i></button>
                                <button type="button" class="btn btn-link text-info p-0 me-2 btn-action-view" title="View Details" data-uuid="${targetId}"><i class="ti ti-eye fs-16"></i></button>
                                <button type="button" class="btn btn-link text-primary p-0 me-2 btn-action-edit" title="Edit" data-uuid="${targetId}"><i class="ti ti-edit fs-16"></i></button>
                                ${statusIcon}
                                <button type="button" class="btn btn-link text-primary p-0 me-2 btn-action-stc" title="Sent To Checker" data-uuid="${targetId}" data-name="${safeName}"><i class="ti ti-arrow-forward-up fs-16"></i></button>
                                <button type="button" class="btn btn-link text-success p-0 btn-action-verify" title="Verify" data-uuid="${targetId}" data-name="${safeName}"><i class="ti ti-check fs-16"></i></button>
                            </div>
                        `;
                    }
                },
                {
                    field: 'name',
                    colId: 'name',
                    headerName: 'Designation Name',
                    width: 240,
                    minWidth: 180,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return `<span class="fw-semibold text-dark"><i class="ti ti-briefcase me-1 text-primary"></i>${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'code',
                    colId: 'code',
                    headerName: 'Code',
                    width: 120,
                    minWidth: 90,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return params.value ? `<span class="badge bg-light text-dark border fs-12">${params.value}</span>` : '<span class="text-muted">-</span>';
                    }
                },
                {
                    field: 'title',
                    colId: 'title',
                    headerName: 'Title',
                    width: 200,
                    minWidth: 150,
                    sortable: true,
                    filter: true
                },
                {
                    field: 'op_status_text',
                    colId: 'actual_status',
                    headerName: 'OP-Status',
                    width: 130,
                    minWidth: 110,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        let text = params.value || 'Normal';
                        let badgeClass = 'bg-primary-transparent text-primary';
                        if (text === 'New Add') badgeClass = 'bg-info-transparent text-info';
                        else if (text === 'Add-Edit') badgeClass = 'bg-warning-transparent text-warning';
                        else if (text === 'Edit') badgeClass = 'bg-secondary-transparent text-secondary';
                        else if (text === 'Delete') badgeClass = 'bg-danger-transparent text-danger';
                        else if (text === 'Deactive') badgeClass = 'bg-dark-transparent text-dark';
                        else if (text === 'Active') badgeClass = 'bg-success-transparent text-success';
                        return `<span class="badge ${badgeClass} fw-semibold">${text}</span>`;
                    }
                },
                {
                    field: 'status_text',
                    colId: 'status',
                    headerName: 'Status',
                    width: 110,
                    minWidth: 90,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        let row = params.data || {};
                        if (row.status === 1) {
                            return `<span class="badge bg-success-transparent text-success fw-semibold"><i class="ti ti-check me-1"></i>Active</span>`;
                        }
                        return `<span class="badge bg-danger-transparent text-danger fw-semibold"><i class="ti ti-x me-1"></i>Deactive</span>`;
                    }
                },
                {
                    field: 'process_status_text',
                    colId: 'process_status',
                    headerName: 'Process Status',
                    width: 140,
                    minWidth: 120,
                    sortable: true,
                    cellRenderer: function (params) {
                        let row = params.data || {};
                        let p = row.process_status;
                        if (p === 1) {
                            return `<span class="badge bg-warning-transparent text-warning fw-semibold"><i class="ti ti-clock me-1"></i>Sent to Checker</span>`;
                        } else if (p === 2) {
                            return `<span class="badge bg-danger-transparent text-danger fw-semibold"><i class="ti ti-arrow-back-up me-1"></i>Returned</span>`;
                        } else if (p === 3) {
                            return `<span class="badge bg-success-transparent text-success fw-semibold"><i class="ti ti-circle-check me-1"></i>Approved</span>`;
                        }
                        return `<span class="badge bg-secondary-transparent text-secondary fw-semibold">Draft</span>`;
                    }
                },
                {
                    field: 'message',
                    colId: 'message',
                    headerName: 'Message',
                    width: 180,
                    minWidth: 140,
                    sortable: true,
                    cellRenderer: function (params) {
                        return params.value ? `<span class="text-truncate d-inline-block" style="max-width: 170px;" title="${params.value}">${params.value}</span>` : '<span class="text-muted">-</span>';
                    }
                },
                {
                    field: 'checker_text',
                    colId: 'checker',
                    headerName: 'Checker',
                    width: 150,
                    minWidth: 120,
                    sortable: true,
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12">${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'created_at',
                    colId: 'created_at',
                    headerName: 'Created At',
                    width: 170,
                    minWidth: 150,
                    sortable: true,
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12"><i class="ti ti-calendar me-1"></i>${params.value || '-'}</span>`;
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

    function openAddModal() {
        $('#form_designation')[0].reset();
        $('#designation_uuid').val('');
        $('#designation_id').val('');
        $('#modal_designation_title').text('Add Designation');
        $('#btn_save_designation').text('Save Designation');
        $('#modal_designation_form').modal('show');
    }

    function openEditModal(targetId) {
        if (!targetId) {
            showNotify('warning', 'Please select a designation to edit');
            return;
        }
        $.get(`/reference/designation/edit/${targetId}`)
            .done(function (res) {
                let data = (res && res.data) || {};
                $('#form_designation')[0].reset();
                $('#designation_uuid').val(data.uuid || targetId);
                $('#designation_id').val(data.id || '');
                $('#designation_name').val(data.name || '');
                $('#designation_code').val(data.code || '');
                $('#designation_title').val(data.title || '');
                $('#designation_details').val(data.details || '');
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
        if (!targetId) {
            showNotify('warning', 'Please select a designation to view');
            return;
        }
        $.get(`/reference/designation/view/${targetId}`)
            .done(function (res) {
                let d = (res && res.data) || {};
                $('#view_name').text(d.name || '-');
                $('#view_code').text(d.code || '-');
                $('#view_title').text(d.title || '-');
                $('#view_details').text(d.details || '-');

                let opText = 'New Add';
                if (d.actual_status === 2) opText = 'Add-Edit';
                else if (d.actual_status === 3) opText = 'Edit';
                else if (d.actual_status === 4) opText = 'Delete';
                else if (d.actual_status === 5) opText = 'Deactive';
                else if (d.actual_status === 6) opText = 'Active';
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
        if (!targetId) {
            showNotify('warning', 'Please select a designation to delete');
            return;
        }
        $('#form_designation_delete')[0].reset();
        $('#delete_target_id').val(targetId);
        $('#delete_designation_name').text(name || 'Selected Item');
        $('#modal_designation_delete').modal('show');
    }

    function openStcModal(targetId, name) {
        if (!targetId) {
            showNotify('warning', 'Please select a designation to send to checker');
            return;
        }
        $('#form_designation_stc')[0].reset();
        $('#stc_target_id').val(targetId);
        $('#stc_designation_name').text(name || 'Selected Item');
        $('#modal_designation_stc').modal('show');
    }

    function openVerifyModal(targetId, name) {
        if (!targetId) {
            showNotify('warning', 'Please select a designation to verify');
            return;
        }
        $('#form_designation_verify')[0].reset();
        $('#verify_target_id').val(targetId);
        $('#verify_designation_name').text(name || 'Selected Item');
        $('#verify_type_approve').prop('checked', true);
        $('#group_return_reason').addClass('d-none');
        $('#modal_designation_verify').modal('show');
    }

    function openStatusModal(actionType, targetId, name) {
        if (!targetId) {
            showNotify('warning', `Please select a designation to ${actionType}`);
            return;
        }
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

    function openHistoryModal(targetId) {
        if (!targetId) return;
        $.get(`/reference/designation/history/${targetId}`)
            .done(function (res) {
                let list = (res && res.data) || [];
                let html = '';
                if (!list || list.length === 0) {
                    html = '<tr><td colspan="9" class="text-center py-3 text-muted">No audit history recorded yet</td></tr>';
                } else {
                    list.forEach(function (h, idx) {
                        let op = 'New Add';
                        if (h.actual_status === 2) op = 'Add-Edit';
                        else if (h.actual_status === 3) op = 'Edit';
                        else if (h.actual_status === 4) op = 'Delete';
                        else if (h.actual_status === 5) op = 'Deactive';
                        else if (h.actual_status === 6) op = 'Active';

                        let ps = h.process_status === 3 ? 'Approved' : (h.process_status === 1 ? 'Sent to Checker' : (h.process_status === 2 ? 'Returned' : 'Draft'));
                        let st = h.status === 1 ? 'Active' : 'Deactive';

                        html += `
                            <tr>
                                <td>${idx + 1}</td>
                                <td>${h.name || '-'}</td>
                                <td>${h.code || '-'}</td>
                                <td>${h.title || '-'}</td>
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

    function bindTopActions() {
        $('#btn_top_add').on('click', function () {
            openAddModal();
        });

        $('#btn_top_delete').on('click', function () {
            if (!selectedDesignation) {
                showNotify('warning', 'Please select a row from the grid first');
                return;
            }
            openDeleteModal(selectedDesignation.uuid || selectedDesignation.id, selectedDesignation.name);
        });

        $('#btn_top_view').on('click', function () {
            if (!selectedDesignation) {
                showNotify('warning', 'Please select a row from the grid first');
                return;
            }
            openViewModal(selectedDesignation.uuid || selectedDesignation.id);
        });

        $('#btn_top_edit').on('click', function () {
            if (!selectedDesignation) {
                showNotify('warning', 'Please select a row from the grid first');
                return;
            }
            openEditModal(selectedDesignation.uuid || selectedDesignation.id);
        });

        $('#btn_top_active').on('click', function () {
            if (!selectedDesignation) {
                showNotify('warning', 'Please select a row from the grid first');
                return;
            }
            openStatusModal('active', selectedDesignation.uuid || selectedDesignation.id, selectedDesignation.name);
        });

        $('#btn_top_deactive').on('click', function () {
            if (!selectedDesignation) {
                showNotify('warning', 'Please select a row from the grid first');
                return;
            }
            openStatusModal('deactivate', selectedDesignation.uuid || selectedDesignation.id, selectedDesignation.name);
        });

        $('#btn_top_stc').on('click', function () {
            if (!selectedDesignation) {
                showNotify('warning', 'Please select a row from the grid first');
                return;
            }
            openStcModal(selectedDesignation.uuid || selectedDesignation.id, selectedDesignation.name);
        });

        $('#btn_top_verify').on('click', function () {
            if (!selectedDesignation) {
                showNotify('warning', 'Please select a row from the grid first');
                return;
            }
            openVerifyModal(selectedDesignation.uuid || selectedDesignation.id, selectedDesignation.name);
        });
    }

    function bindRowActions() {
        $(document).on('click', '.btn-action-delete', function (e) {
            e.stopPropagation();
            let uuid = $(this).data('uuid');
            let name = decodeURIComponent($(this).data('name') || '');
            openDeleteModal(uuid, name);
        });

        $(document).on('click', '.btn-action-view', function (e) {
            e.stopPropagation();
            let uuid = $(this).data('uuid');
            openViewModal(uuid);
        });

        $(document).on('click', '.btn-action-edit', function (e) {
            e.stopPropagation();
            let uuid = $(this).data('uuid');
            openEditModal(uuid);
        });

        $(document).on('click', '.btn-action-active', function (e) {
            e.stopPropagation();
            let uuid = $(this).data('uuid');
            let name = decodeURIComponent($(this).data('name') || '');
            openStatusModal('active', uuid, name);
        });

        $(document).on('click', '.btn-action-deactive', function (e) {
            e.stopPropagation();
            let uuid = $(this).data('uuid');
            let name = decodeURIComponent($(this).data('name') || '');
            openStatusModal('deactivate', uuid, name);
        });

        $(document).on('click', '.btn-action-stc', function (e) {
            e.stopPropagation();
            let uuid = $(this).data('uuid');
            let name = decodeURIComponent($(this).data('name') || '');
            openStcModal(uuid, name);
        });

        $(document).on('click', '.btn-action-verify', function (e) {
            e.stopPropagation();
            let uuid = $(this).data('uuid');
            let name = decodeURIComponent($(this).data('name') || '');
            openVerifyModal(uuid, name);
        });

        $('#btn_open_history_from_view').on('click', function () {
            let uuid = $(this).data('uuid');
            $('#modal_designation_view').modal('hide');
            openHistoryModal(uuid);
        });

        $('input[name="verify_type"]').on('change', function () {
            if ($(this).val() === '0') {
                $('#group_return_reason').removeClass('d-none');
                $('#verify_return_reason').prop('required', true);
            } else {
                $('#group_return_reason').addClass('d-none');
                $('#verify_return_reason').prop('required', false);
            }
        });
    }

    function bindFormSubmissions() {
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
                    showNotify('danger', msg);
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
