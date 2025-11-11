document.addEventListener("DOMContentLoaded", function() {
    const newDictCreated = {};
    const dictsToCreate  = [];

    let manualAvitoIds = "";
    let currentManualCtx = null;

    const ALLOW_NEW_DICTS = new Set([
        'reklId',
        'buyingTypeId',
        'profId',
        'providerId',
        'trafficSourceId'
    ]);

    const CPA_REQUIRED_EXPORT_TYPE = 9;

    function enableSelect2(scope = document, tryAllowNew = false) {
        const $scope   = $(scope);
        const $targets = $scope.is('select')
            ? $scope.filter(':not(.select2-hidden-accessible)')
            : $scope.find('select:not(.select2-hidden-accessible)');

        $targets.each(function () {
            const multiple    = this.multiple;
            const placeholder = $(this).find('option[disabled]:first').text() || 'Выберите…';
            const dictName    = this.name || '';
            const allowNew = tryAllowNew && ALLOW_NEW_DICTS.has(dictName);

            $(this).select2({
                width:        '100%',
                language:     'ru',
                placeholder,
                closeOnSelect: !multiple,
                allowClear:   !multiple,
                tags: allowNew,
                createTag: allowNew ? params => {
                    const term = $.trim(params.term);
                    if (!term) return null;

                    const randNegId = -Math.floor(Math.random() * 1000);

                    return { id: randNegId, text: term, newTag: true };
                } : undefined,
                templateResult: data => {
                    if (data.newTag) {
                        return $('<span>')
                            .append('<span class="me-1 text-success">+</span>')
                            .append(document.createTextNode(data.text));
                    }
                    return data.text;
                },
                insertTag: allowNew ? (data, tag) => { data.push(tag); } : undefined
            });

            if (allowNew) {
                $(this).on('select2:select', function (e) {
                    const dict = this.name;
                    const item = e.params.data;

                    if (item.newTag && !newDictCreated[dict]) {
                        newDictCreated[dict] = true;
                        dictsToCreate.push({ dict, id: item.id, name: item.text });
                    }
                });
            }
        });
    }

    enableSelect2(document, true);

    $('#navVacancyAdmin').addClass('active');

    function toInt(v) {
        return parseInt(v || "0", 10) || 0;
    }

    function hasOffers() {
        return saveState.offers.length > 0;
    }

    function addClass(id, className) {
        const el = document.getElementById(id);
        if (el && !el.classList.contains(className)) {
            el.classList.add(className);
        }
    }

    function deleteClass(id, className) {
        const el = document.getElementById(id);
        if (!el) return;
        if (el && el.classList.contains(className)) {
            el.classList.remove(className);
        }
    }

    function fetchJSON(url) {
        return fetch(url).then(r => {
            if (!r.ok) return Promise.reject(new Error(`Network error: ${r.status}`));
            return r.json();
        });
    }

    let allFilters = [];
    let questionTypesList = [];

    async function loadAllFilters() {
        try {
            const data = await fetchJSON('/rk/filters');
            allFilters = data.map(f => ({
                id:   f.Id   ?? f.id,
                name: f.Name ?? f.name,
                value:f.Value?? f.value
            }));
        } catch (e) {
            console.error('Не удалось загрузить список фильтров', e);
        }
    }

    function rebuildFilterSelect(offerId) {
        const $select = $(`#newFilterRow-${offerId} select[name="filterId"]`);
        if (!$select.length) return;

        const offer = saveState.offers.find(o =>
            o.id === offerId || o.tempId === offerId
        );

        const existingIds = new Set(offer ? offer.filter_ids : []);

        $select.empty();

        allFilters.forEach(f => {
            if (!existingIds.has(f.id)) {
                const opt = new Option(`${f.name} (${f.value})`, f.id);
                opt.dataset.name = f.name;
                opt.dataset.val  = f.value;
                $select.append(opt);
            }
        });

        if ($select.hasClass('select2-hidden-accessible')) {
            $select.trigger('change.select2');
        } else {
            enableSelect2($select);
        }

        if (!$select.children().length) {
            $select
                .append(new Option("Нет доступных фильтров", "", true, false))
                .prop("disabled", true);
        } else {
            $select.prop("disabled", false);
        }
    }

    function initModalHandlers() {
        $('#deleteOffer').on('show.bs.modal', function(event) {
            const button = $(event.relatedTarget);
            const vac2vacId = button.data('bs-whatever');
            $(this).find('input[name="vac2vacId"]').val(vac2vacId);
        });

        $('#deleteFilter').on('show.bs.modal', function(event) {
            const button = $(event.relatedTarget);
            const offerSqlFilterId = button.data('bs-whatever');
            $(this).find('input[name="offerSqlFilterId"]').val(offerSqlFilterId);
        });
    }

    function buildDictSelect(id, list, current) {
        const opts = list.map(it => {
            const value = it.Id   ?? it.id;
            const label = it.Name ?? it.name;
            return `<option value="${value}" ${value === current ? 'selected' : ''}>${label}</option>`;
        }).join('');

        return `<select id="${id}" name="${id}" class="form-select dict-select">${opts}</select>`;
    }

    function renderVacancyDicts() {
        const vTxt   = document.getElementById('vacancy-text');
        if (!vTxt) return;

        const cur = {
            reklId:           toInt(vTxt.dataset.reklId),
            buyingTypeId:     toInt(vTxt.dataset.buyingTypeId),
            profId:           toInt(vTxt.dataset.profId),
            providerId:       toInt(vTxt.dataset.providerId),
            trafficSourceId:  toInt(vTxt.dataset.trafficSourceId),
            exportTypeId:     toInt(vTxt.dataset.exportTypeId)
        };

        const info = {};
        (vTxt.textContent || '')
            .split('\n')
            .map(l => l.trim())
            .filter(l => l && l.includes(':'))
            .forEach(l => {
                const [k, ...rest] = l.split(':');
                info[k.trim()] = rest.join(':').trim();
            });

        const buildText = (id, val = '') => {
            const safeVal = escapeHtmlAttr(val);
            return `<input id="${id}" name="${id}" class="form-control" value="${safeVal}">`;
        };

        const html = `
          <table class="table table-bordered">
            <thead>
              <tr>
                <th>ID РК</th>
                <th>Рекламодатель</th>
                <th>Профессия</th>
                <th>Провайдер</th>
                <th>Тип передачи</th>
                <th>Тип покупки</th>
                <th>Источник трафика</th>
                <th>Внешний ID</th>
                <th>Прочее</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>${info['ИД РК']            || ''}</td>
                <td>${buildDictSelect('reklId',           reklsList,          cur.reklId)}</td>
                <td>${buildDictSelect('profId',           professionsList,    cur.profId)}</td>
                <td>${buildDictSelect('providerId',       providersList,      cur.providerId)}</td>
                <td>${buildDictSelect('exportTypeId',     exportTypesList,    cur.exportTypeId)}</td>
                <td>${buildDictSelect('buyingTypeId',     buyingTypesList,    cur.buyingTypeId)}</td>
                <td>${buildDictSelect('trafficSourceId',  trafficSourcesList, cur.trafficSourceId)}</td>
                <td>${buildText('cpaExId',           vTxt.dataset.cpaExId           || '')}</td>
                <td>${buildText('otherSpecification', vTxt.dataset.otherSpecification || '')}</td>
              </tr>
            </tbody>
          </table>`
        ;

        $('#vacancy-dicts').html(html);
        $('#vacancy-openai-card #vacancy-dicts').html(html);
        enableSelect2($('#vacancy-card #vacancy-dicts select:visible'), true);
    }

    function initVacancyFormToggle() {
        const checkbox = document.getElementById('open-ai-support');
        function toggleForms() {
            const isOpenAI = checkbox.checked;

            if (isOpenAI) {
                addClass('vacancy-card', 'd-none');
                deleteClass('vacancy-openai-card', 'd-none');
                enableSelect2($('#vacancy-openai-card #vacancy-dicts select'), true);
            } else {
                addClass('vacancy-openai-card', 'd-none');
                deleteClass('vacancy-card', 'd-none');
                enableSelect2($('#vacancy-card #vacancy-dicts select'), true);
            }
        }

        checkbox.addEventListener('change', toggleForms);
        toggleForms();
    }

    function initExportTypeSelects() {
        $('.export-type-select').each(function() {
            const $select = $(this);
            const offerId = $select.data('offer-id');
            const extraFields = $(`#exportExtraFields-${offerId}`);

            function onTypeChange() {
                const selectedVal = $select.val();
                if (!selectedVal) {
                    extraFields.hide().html("");
                    return;
                }

                if (selectedVal === "2") {
                    const sheetId = $select.data('client-sheet-id') || "";
                    const tabName = $select.data('client-tab-name') || "";
                    extraFields.html(
                        '<div class="row">' +
                        '<div class="col-md-6">' +
                        `<label for="sheetId-${offerId}" class="form-label fw-bold">Sheet ID</label>` +
                        `<input type="text" class="form-control" id="sheetId-${offerId}" name="sheetId" value="${sheetId}">` +
                        '</div>' +
                        '<div class="col-md-6">' +
                        `<label for="tabName-${offerId}" class="form-label fw-bold">Название вкладки</label>` +
                        `<input type="text" class="form-control" id="tabName-${offerId}" name="tabName" value="${tabName}">` +
                        '</div>' +
                        '</div>'
                    );
                    extraFields.show();
                } else {
                    const link = $select.data('client-link') || "";
                    extraFields.html(
                        '<div class="row">' +
                        '<div class="col-md-12">' +
                        `<label for="link-${offerId}" class="form-label fw-bold">Ссылка</label>` +
                        `<input type="text" class="form-control" id="link-${offerId}" name="link" value="${link}">` +
                        '</div>' +
                        '</div>'
                    );
                    extraFields.show();
                }
            }

            $select.on('change.select2', onTypeChange);
            onTypeChange();
        });
        enableSelect2('.export-type-select');
    }

    function initOffersCollapse() {
        const offers = $('[id^="editOffer-"]');
        if (offers.length === 1) {
            offers.first().collapse('show');
        } else {
            offers.collapse('hide');
        }
    }

    const saveState = {
        vacancy: null,
        offers: [],
        questionnaire: null
    };

    function gatherVacancy() {
        const isOpenAI = $('#open-ai-support').is(':checked');
        const $form    = isOpenAI ? $('#vacancy-openai-card') : $('#vacancy-card');

        const val = sel => $form.find(sel).val();
        const chk = sel => $form.find(sel).is(':checked');

        const base = {
            id: toInt($('#vacancyId').val()),
            loadToZP: chk('[name="loadToZP"]'),
            questionnaire: toInt($('#questId').val()),
            description: val('#description') || '',
            helloText: val('#helloText') || '',
            byeText:   val('#byeText') || '',
            enablePostProcessing: chk('#enablePostProcessing'),
            openAiSupport: isOpenAI,
            postProcessingMessages: chk('#enablePostProcessing') ? val('#postProcessingMessages') : '',
            dialogLifeTimeInMinutes:            toInt(val('#' + (isOpenAI ? 'openai-' : '') + 'dialogLifeTimeInMinutes')),
            followUpMessageIntervalInMinutes:   toInt(val('#' + (isOpenAI ? 'openai-' : '') + 'followUpMessageIntervalInMinutes')),
            countOfMessagesAfterFinishedDialog: toInt(val('#' + (isOpenAI ? 'openai-' : '') + 'countOfMessagesAfterFinishedDialog')),
            reklId: toInt(val('#reklId')),
            buyingTypeId: toInt(val('#buyingTypeId')),
            profId: toInt(val('#profId')),
            providerId: toInt(val('#providerId')),
            trafficSourceId: toInt(val('#trafficSourceId')),
            exportTypeId: toInt(val('#exportTypeId')),
            cpaExId:            (val('#cpaExId') || '').trim(),
            otherSpecification: (val('#otherSpecification') || '').trim(),
        };

        if (!isOpenAI) return base;

        return {
            ...base,
            openAIDetails: {
                id: parseInt(val('#openAIid')) || 0,
                vacancyDescription:   val('#vacancyDescription'),
                assistantDescription: val('#assistantDescription'),
                questions:            val('#questions'),
                assistantTemperature: parseFloat(val('#assistantTemperature')) || 0,
            }
        };
    }

    function collectOfferFromDOM(domId) {
        const $row = $(`tr[data-temp-id='${domId}'], tr[data-offer-id='${domId}'], tr[data-v2o-id='${domId}']`);
        if (!$row.length) {
            return { id: 0, name: "", sql_for_zp: "", zp_script_name: "", enable: false,
                ignore_name: false, export_zp_pg: false, export_type_id: 0, client_id: 0,
                link: "", sheet_id: "", sheet_name: "", filter_ids: [] };
        }

        const offerObj = saveState.offers.find(o =>
            toInt(o.id) === domId || toInt(o.tempId) === domId
        );

        let realId = 0;
        if (offerObj) {
            realId = toInt(offerObj.id);
        } else {
            realId = toInt($row.find('[name="offerId"]').val());
        }

        return {
            id: realId,
            name: $row.find('[name="offerName"]').val(),
            sqlForZp:      $row.find('[name="sqlForZp"]').val(),
            zpScriptName:  $row.find('[name="zpScriptName"]').val(),
            enable: $row.find('[name="enable"]').is(':checked'),
            ignoreName:    $row.find('[name="ignoreName"]').is(':checked'),
            exportZpPg:    $row.find('[name="exportZpPg"]').is(':checked'),
            exportTypeId:  toInt($row.find('[name="exportTypeId"]').val()),
            link: $row.find('[name="link"]').val() || "",
            sheetId: $row.find('[name="sheetId"]').val() || "",
            sheetName: $row.find('[name="tabName"]').val() || "",
            filter_ids: $row.find('tr[data-filter-id]').map(function() {
                return parseInt(this.dataset.filterId, 10);
            }).get().filter(id => Number.isFinite(id))
        };
    }

    function gatherQuestionnaire() {
        const $rows = $('#questionnaire-block tbody tr')
            .not('#new_tr')
            .filter(function () {
                const $r = $(this);
                return !$r.is('[data-deleted="1"]') && !$r.is('[data-template]');
            });

        const qId   = toInt($('#questId').val());
        const qName = '';

        const questions = $rows.map(function (idx) {
            const $row = $(this);

            let id = toInt($row.find('input[name="questionid"]').val());

            if (!id) {
                id = toInt($row.attr('data-question-id'));
            }

            if (!id || id < 0 ) id = 0;

            return {
                id:         id,
                text:       ($row.find('textarea[id^="questonName"]').val()        || '').trim(),
                wrongAnswer:($row.find('textarea[id^="wrongAnswerMessages"]').val()|| '').trim(),
                followUp:   ($row.find('textarea[id^="followUpMessages"]').val()   || '').trim(),
                typeId:     toInt($row.find('select.question-type-select').val()),
                isRequired: $row.find('input[id^="isRequired"]').is(':checked'),
                sort:       idx + 1
            };
        }).get();

        return { id: qId, name: qName, questions };
    }

    async function renderNewOfferUI(offer) {
        const domId = offer.tempId;
        const filterOptions =
            allFilterOptions ||
            '<option value="" disabled>Нет доступных фильтров</option>';

        const listRow = `
        <tr data-v2o-id="${domId}" data-temp-id="${domId}" class="table-success">
            <td style="width:10px;">${offer.name || '(без названия)'}</td>
            <td class="text-start" style="width:150px;">
                <button
                    class="btn btn-secondary"
                    type="button"
                    data-bs-toggle="collapse"
                    data-bs-target="#editOffer-${domId}"
                    aria-expanded="false"
                    aria-controls="editOffer-${domId}">
                    Редактировать
                </button>
            </td>
            <td class="text-end" style="width:100px;">
                <button
                    class="btn btn-primary"
                    data-bs-toggle="modal"
                    data-bs-target="#deleteOffer"
                    data-bs-whatever="${domId}">
                    Удалить
                </button>
            </td>
        </tr>`;

        const editRow = `
        <tr class="collapse" id="editOffer-${domId}" data-v2o-id="${domId}" data-offer-id="${domId}" data-temp-id="${domId}">
            <td colspan="3" class="p-0">
                <form class="row g-3" data-temp="1">
                    <input type="hidden" name="offerId" value="0">
                    <input type="hidden" name="vacancyId" value="${$('#vacancyId').val()}">

                    <div class="col-md-2">
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="enable"
                                ${offer.enable ? 'checked' : ''}>
                            <label class="form-check-label fw-bold">Активен</label>
                        </div>
                    </div>
                    <div class="col-md-3">
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="ignoreName"
                                ${offer.ignoreName ? 'checked' : ''}>
                            <label class="form-check-label fw-bold">Игнорировать имя лида</label>
                        </div>
                    </div>
                    <div class="col-md-3">
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" name="exportZpPg"
                                ${offer.exportZpPg ? 'checked' : ''}>
                            <label class="form-check-label fw-bold">Экспорт ZP PG</label>
                        </div>
                    </div>

                    <div class="col-md-6">
                        <label class="form-label fw-bold">Наименование</label>
                        <input type="text" class="form-control" name="offerName" value="${offer.name}" readonly>
                    </div>
                    <div class="col-md-6">
                        <label class="form-label fw-bold">Имя скрипта ZP</label>
                        <input type="text" class="form-control" name="zpScriptName"
                            value="${offer.zpScriptName}">
                    </div>

                    <div class="col-md-12">
                        <label class="form-label fw-bold">SQL запрос</label>
                        <textarea class="form-control" rows="3" name="sqlForZp"
                            style="white-space: pre-wrap;">${offer.sqlForZp}</textarea>
                    </div>

                    <div class="col-md-6">
                        <label class="form-label fw-bold">Тип экспорта</label>
                        <select
                            name="exportTypeId"
                            class="form-select export-type-select"
                            data-offer-id="${domId}"
                            data-current="${offer.exportTypeId}"
                            data-client-sheet-id="${offer.sheetId}"
                            data-client-tab-name="${offer.sheetName}"
                            data-client-link="${offer.link}">
                            ${$('#export_type_id').html()}
                        </select>

                        <div id="exportExtraFields-${domId}" class="mt-3" style="display: none;"></div>
                    </div>

                    <div class="col-md-12">
                        <button type="button" class="btn btn-primary d-none btn-save-individual">Сохранить</button>
                    </div>
                </form>

                <div class="mt-3 px-0">
                    <h5>Фильтры</h5>
                    <div class="table-responsive">
                        <table class="table table-default table-sm filters-table">
                        <colgroup>
                            <col style="width:35%">
                            <col style="width:50%">
                            <col style="width:10%">
                        </colgroup>
                            <thead>
                                <tr><th>Наименование</th><th class="text-start">Значение</th><th></th></tr>
                            </thead>
                            <tbody>
                                <tr id="newFilterRow-${domId}" class="d-none align-middle">
                                    <form class="add-filter-form" method="post" action="/dict/offers/add_filter">
                                      <input type="hidden" name="offerId" value="${domId}">
                                      <td colspan="2">
                                        <select class="form-select filter-select" name="filterId" style="max-width:300px;">
                                          ${filterOptions}
                                        </select>
                                      </td>
                                      <td class="text-end">
                                        <button type="submit" class="btn btn-primary">Добавить</button>
                                      </td>
                                    </form>
                                </tr>
                                ${offer.filter_ids.map(fid => {
            return `<tr data-filter-id="${fid}" class="table-success">
                                        <td>–</td>
                                        <td>–</td>
                                        <td class="text-end">
                                            <button
                                                class="btn btn-primary"
                                                data-bs-toggle="modal"
                                                data-bs-target="#deleteFilter"
                                                data-bs-whatever="${fid}">
                                                Удалить
                                            </button>
                                        </td>
                                    </tr>`;
        }).join('')}
                            </tbody>
                            <tfoot>
                                <tr>
                                    <td colspan="3">
                                        <button
                                            type="button"
                                            onclick="deleteClass('newFilterRow-${domId}', 'd-none')"
                                            class="btn btn-secondary">
                                            Добавить фильтр
                                        </button>
                                    </td>
                                </tr>
                            </tfoot>
                        </table>
                    </div>
                </div>
            </td>
        </tr>`;

        $('#newOfferForm').before(listRow + editRow);
        enableSelect2(`#editOffer-${domId}`);
        $(`#editOffer-${domId}`).collapse({ toggle: false });

        await refreshFilterOptions(domId);

        populateExportTypeSelect(domId, offer.exportTypeId, offer.sheetId, offer.sheetName, offer.link);
    }

    function populateExportTypeSelect(offerId, currentExportTypeId, sheetId, sheetName, link) {
        const $sel = $(`select[name="exportTypeId"][data-offer-id="${offerId}"]`);
        if (!$sel.length || !Array.isArray(exportTypesList)) return;

        $sel.empty();
        $sel.append('<option value="">Выбрать тип экспорта</option>');
        exportTypesList.forEach(et => {
            $sel.append(`<option value="${et.Id}">${et.Name}</option>`);
        });
        if (currentExportTypeId) {
            $sel.val(currentExportTypeId);
        }

        const extraFields = $(`#exportExtraFields-${offerId}`);

        function onTypeChange() {
            const selectedVal = $sel.val();
            if (!selectedVal) {
                extraFields.hide().html("");
                return;
            }

            if (selectedVal === "2") {
                extraFields.html(
                    '<div class="row">' +
                    '<div class="col-md-6">' +
                    `<label for="sheetId-${offerId}" class="form-label fw-bold">Sheet ID</label>` +
                    `<input type="text" class="form-control" id="sheetId-${offerId}" name="sheetId" value="${sheetId || ""}">` +
                    '</div>' +
                    '<div class="col-md-6">' +
                    `<label for="tabName-${offerId}" class="form-label fw-bold">Название вкладки</label>` +
                    `<input type="text" class="form-control" id="tabName-${offerId}" name="tabName" value="${sheetName || ""}">` +
                    '</div>' +
                    '</div>'
                );
                extraFields.show();
            } else {
                extraFields.html(
                    '<div class="row">' +
                    '<div class="col-md-12">' +
                    `<label for="link-${offerId}" class="form-label fw-bold">Ссылка</label>` +
                    `<input type="text" class="form-control" id="link-${offerId}" name="link" value="${link || ""}">` +
                    '</div>' +
                    '</div>'
                );
                extraFields.show();
            }
        }

        $sel.off('change', onTypeChange).on('change', onTypeChange);
        onTypeChange();
    }

    function resetSaveState() {
        saveState.vacancy = null;
        saveState.offers = [];
    }

    function showError(msg) {
        let $alert = $('#dynamicErrorAlert');
        if (!$alert.length) {
            $alert = $('<div id="dynamicErrorAlert" class="alert alert-danger" role="alert"></div>')
                .prependTo('.main');
        }
        $alert.text(msg);
    }

    async function refreshFilterOptions(offerId) {
        const $select = $(`#newFilterRow-${offerId} select[name="filterId"]`);
        if (!$select.length) return;

        const offer = saveState.offers.find(o =>
            o.id === offerId || o.tempId === offerId
        );

        const existingIds = new Set(offer ? offer.filter_ids : []);

        $select.empty();

        allFilters.forEach(f => {
            if (!existingIds.has(f.id)) {
                const opt = new Option(`${f.name} (${f.value})`, f.id);
                opt.dataset.name = f.name;
                opt.dataset.val  = f.value;
                $select.append(opt);
            }
        });

        if (!$select.children().length) {
            $select
                .append(new Option("Нет доступных фильтров", "", true, false))
                .prop("disabled", true);
        } else {
            $select.prop("disabled", false);
        }
    }

    async function updateOfferFilters(action, offerId, filterId) {
        const offer = saveState.offers.find(o =>
            o.id === offerId || o.tempId === offerId
        );

        if (!offer) return;

        if (action === "add") {
            if (!offer.filter_ids.includes(filterId)) {
                offer.filter_ids.push(filterId);
            }
        } else {
            offer.filter_ids = offer.filter_ids.filter(id => id !== filterId);
        }

        $(`#editOffer-${offerId}`).attr("data-upd", "1");
        await refreshFilterOptions(offerId);
    }

    function updateAddOfferAvailability() {
        if (hasOffers()) {
            $('#newOfferForm').addClass('d-none');
            $('#add-offer-footer').remove();
        } else {
            $('#newOfferForm').addClass('d-none');
            if (!$('#add-offer-footer').length) {
                const btnRow = `
                <tfoot id="add-offer-footer">
                  <tr>
                    <td colspan="3">
                      <button type="button" id="btn-add-offer" class="btn btn-secondary">Добавить оффер</button>
                    </td>
                  </tr>
                </tfoot>`;
                $('table.table-default').append(btnRow);
                $(document).on('click', '#btn-add-offer', () => {
                    $('#newOfferForm').removeClass('d-none');
                    $('#btn-add-offer').hide();
                });
            } else {
                $('#btn-add-offer').show();
            }
        }
    }

    function collectExistingOffers() {
        saveState.offers = $('[id^="editOffer-"]').map(function () {
            return collectOfferFromDOM(toInt(this.dataset.offerId));
        }).get();
    }

    function initHandlers() {
        collectExistingOffers();
        initModalHandlers();
        initVacancyFormToggle();
        initExportTypeSelects();
        initOffersCollapse();

        const $fileInput       = $('#avitoFile');
        const $fileLabel       = $('#avitoFileLabel');

        const $fileInputOpenAI = $('#openaiAvitoFile');
        const $fileLabelOpenAI = $('#openaiAvitoFileLabel');

            function bindFileLabel($input, $label, manualBtnId, infoDivSel) {
            $input.data('manual-btn', manualBtnId);
            $input.data('info-div', infoDivSel);

            const $text = $label.find('.label-text');

            $input.on('change', function () {
                const manualBtn = $('#' + $(this).data('manual-btn'));

                if (this.files && this.files.length) {
                    $label.removeClass('btn-outline-secondary').addClass('btn-success');
                    $text.text(`Файл: ${this.files[0].name}`);

                    manualBtn.prop('disabled', true);

                    const vacId = parseInt($('#vacancyId').val() || '0', 10);
                    const fd = new FormData();
                    fd.append('avitoIdsFile', this.files[0]);
                    fd.append('inputType', 'file');

                    fetch(`/rk/check_avito_ids?vacancyId=${vacId}`, {
                        method: 'POST',
                        body: fd
                    })
                        .then(r => {
                            if (!r.ok) {
                                return r.text().then(t => { throw new Error(t || 'Ошибка проверки Avito-ID'); });
                            }
                            return r.json();
                        })
                        .then(({ new: newCnt, count, total }) => {
                            const qty  = newCnt ?? count ?? total ?? 0;
                            const $info = $($(this).data('info-div'));
                            console.log(qty);
                            console.log($info);
                            if ($info.length) {
                                $info.removeClass('d-none')
                                    .text(`Загружено ${qty} новых ID (из файла)`);
                            }

                            $('#dynamicErrorAlert').remove();
                        })
                        .catch(err => {
                            showError(err.message);
                        });
                } else {
                    $label.removeClass('btn-success').addClass('btn-outline-secondary');
                    $text.text('Загрузить из файла (.xlsx)');

                    if (!manualAvitoIds) {
                        manualBtn.prop('disabled', false);
                    }

                    const $info = $(infoDivSel);
                    if ($info.length) $info.addClass('d-none').text('');
                }
            });
        }

        bindFileLabel($fileInput, $fileLabel, 'btn-manual-avito', '#manualAvitoInfo');
        bindFileLabel($fileInputOpenAI, $fileLabelOpenAI, 'btn-manual-avito-openai', '#manualAvitoInfoOpenAI');

        updateAddOfferAvailability();

        $(document).on('click', '#confirmDeleteOffer', function () {
            const id = toInt($('#vac2vacId').val());
            if (id) {
                saveState.offers = saveState.offers.filter(o => o.id !== id && o.tempId !== id);
                $(`[data-v2o-id='${id}'], [data-offer-id='${id}']`).remove();
                $(`#editOffer-${id}`).remove();
                updateAddOfferAvailability();
            }
            $('#deleteOffer').modal('hide');
        });

        $(document).on('click', '#confirmDeleteFilter', async function() {
            const joinId = toInt($('#offerSqlFilterId').val());
            if (!joinId) {
                $('#deleteFilter').modal('hide');
                return;
            }
            const $row = $(`tr[data-filter-id='${joinId}']`);
            if (!$row.length) {
                $('#deleteFilter').modal('hide');
                return;
            }

            const filterId = parseInt($row.data('filter-id'), 10);
            const $editRow = $row.closest("tr[id^='editOffer-']");
            const offerId = toInt($editRow.data("offer-id"));
            await updateOfferFilters('del', offerId, filterId);

            $('#deleteFilter').modal('hide');
            setTimeout(() => {
                $row.remove();
                rebuildFilterSelect(offerId);
            }, 200);
        });

        $(document).on('submit', 'form[action="/dict/offers/add_filter"]', async function(e) {
            e.preventDefault();

            const $form = $(this);
            const data = $form.serializeArray().reduce((acc, { name, value }) => {
                acc[name] = value;
                return acc;
            }, {});

            const offerId  = toInt(data.offerId);
            const filterId = toInt(data.filterId);

            if (!offerId || !filterId) {
                return;
            }

            const selectEl = this.elements.filterId;
            const selected = selectEl.options[selectEl.selectedIndex];
            if (!selected) return;

            const rawName = selected.dataset.name ?? selected.textContent.split(' (')[0];
            const rawVal  = selected.dataset.val  ?? selected.textContent.split(' (')[1]?.replace(/\)$/, '');

            const name = rawName.trim() || '—';
            const val  = rawVal.trim()  || '—';

            await updateOfferFilters('add', offerId, filterId);

            const rowHtml = `
              <tr data-filter-id="${filterId}" class="table-success">
                <td>${name}</td>
                <td>${val}</td>
                <td>
                  <button
                    class="btn btn-primary"
                    data-bs-toggle="modal"
                    data-bs-target="#deleteFilter"
                    data-bs-whatever="${filterId}">
                    Удалить
                  </button>
                </td>
              </tr>`;

            $(`#newFilterRow-${offerId}`).before(rowHtml);

            rebuildFilterSelect(offerId);

            $form.closest('tr').addClass('d-none');
            $form[0].reset();
        });

        $(document).on('submit', 'form[action="/dict/offers/create"]', async function(e) {
            e.preventDefault();
            const newOffer = {
                id: 0,
                tempId: Date.now(),
                name: $('#newOfferName').val() || "",
                sqlForZp:     $('#newSqlForZP').val() || "",
                zpScriptName: $('#newZpScriptName').val() || "",
                enable: $('#newEnable').is(':checked'),
                ignoreName:   $('#newIgnoreName').is(':checked'),
                exportZpPg:   $('#newExportZPPG').is(':checked'),
                exportTypeId: toInt($('#export_type_id-new').val()),
                link: $('#exportExtraFields-new input[name="link"]').val() || "",
                sheetId: $('#exportExtraFields-new input[name="sheetId"]').val() || "",
                sheetName: $('#exportExtraFields-new input[name="tabName"]').val() || "",
                filter_ids: []
            };
            saveState.offers.push(newOffer);
            await renderNewOfferUI(newOffer);
            this.reset();
            updateAddOfferAvailability();
        });

        $('form[action="/dict/vacancy-openai/update"]').on('submit', e => e.preventDefault());

        $(document).on('change input', '[id^="editOffer-"] :input', function() {
            $(this).closest('[id^="editOffer-"]').attr('data-upd', "1");
        });

        $(document).on('click', '#saveAllBtn', function (e) {
            e.preventDefault();

            const vacId   = parseInt($('#vacancyId').val() || '0', 10);

            saveState.vacancy = gatherVacancy();

            saveState.offers = $('[id^="editOffer-"]').map(function() {
                const domId = toInt(this.dataset.offerId);
                return collectOfferFromDOM(domId);
            }).get();

            if (!saveState.vacancy.openAiSupport) {
                saveState.questionnaire = gatherQuestionnaire();
            }

            if (dictsToCreate.length) {
                saveState.newDictionaries = dictsToCreate;
            }

            if (manualAvitoIds) {
                saveState.manualAvitoIds = manualAvitoIds
            }

            const fd = new FormData();
            fd.append('payload', JSON.stringify(saveState));

            ['avitoFile', 'openaiAvitoFile'].forEach(id => {
                const fi = document.getElementById(id);
                if (fi && fi.files && fi.files.length) {
                    fd.append('avitoIdsFile', fi.files[0]);
                }
            });

            fetch('/dict/vacancy_admin/save_all', { method: 'POST', body: fd })
                .then(r => {
                    if (!r.ok) return r.text().then(t => { throw new Error(t || 'Ошибка сохранения'); });
                    return r.json();
                })
                .then(data => {
                    resetSaveState();
                    if (vacId === 0 && data.vacancyId) {
                        window.location.href = `/dict/vacancy_admin/edit/${data.vacancyId}?afterSave=1`;
                    } else {
                        reloadPageAfterSaving();
                    }
                })
                .catch(err => showError(err.message));
        });

        $(document).on('click', 'button#up, button#down', function (e) {
            e.preventDefault();
            const $row = $(this).closest('tr');
            if (this.id === 'up') {
                $row.prevAll('tr').not('#new_tr').first().before($row);
            } else {
                $row.nextAll('tr').not('#new_tr').first().after($row);
            }
        });

        $(document).off('submit', 'form[action^="/dict/questions"]');

        $(document).off('click.saveQuestion');

        $(document).on('click.saveQuestion', 'button[id^="saveQuestion"]', function (e) {
            e.preventDefault();

            const $btn = $(this);
            const $row = $btn.closest('tr[data-question-id]');

            $row.find('textarea').prop('readonly', true);
            $row.find('select').prop('disabled', true);
            $row.find('input[type="checkbox"]').prop('disabled', true);

            $btn.prop('disabled', true);

            $row.find('button.btn-secondary').prop('disabled', false);

            $row.attr('data-upd', '1');
        });

        $(document).off('click', '.add-question-btn');

        $(document).on('click', '.add-question-btn', function () {
            const $btn    = $(this);
            const $form   = $btn.closest('form');
            const $tplRow = $btn.closest('tr');
            if (!$tplRow.length) return;

            const domKey = Date.now() * -1;

            const $newTpl = $tplRow.clone(false);

            if ($newTpl.find('form').length === 0) {
                const $cells = $newTpl.children().detach();
                $('<form>', {
                    method : 'post',
                    class  : 'row g-3',
                    action : '/dict/questions/add'
                }).append($cells).appendTo($newTpl);
            }

            $tplRow
                .attr('data-question-id', domKey)
                .removeAttr('id data-template')
                .removeClass('d-none')
                .addClass('table-success');

            $tplRow.find('textarea,select').prop({ readonly: true, disabled: true });
            $tplRow.find('input[type="checkbox"]').prop('disabled', true);

            $btn.prop('disabled', true)
                .text('Добавлено')
                .removeClass('btn-primary')
                .addClass('btn-secondary');

            $tplRow.find('[id]').each(function () {
                const id = this.id;
                if (/_?New$/.test(id)) this.id = id.replace(/_?New$/, domKey);
            });
            $tplRow.find('[name]').each(function () {
                const nm = this.name;
                if (/_?New$/.test(nm)) this.name = nm.replace(/_?New$/, domKey);
            });

            const $hiddenId = $tplRow.find('input[name="questionid"]');
            if ($hiddenId.length) {
                $hiddenId.val('0');
            } else {
                const $rowForm = $tplRow.find('form').first();
                $('<input>', { type: 'hidden', name: 'questionid', value: '0' }).prependTo($rowForm.length ? $rowForm : $tplRow);
            }

            $newTpl
                .attr({ id: 'new_tr', 'data-template': '1' })
                .addClass('d-none')
                .removeClass('table-success')
                .find('textarea').val('').prop({ readonly: false }).end()
                .find('select').val('').prop({ disabled: false }).end()
                .find('input[type="checkbox"]').prop({ checked: false, disabled: false }).end()
                .find('input[name="questionid"]').remove();

            $newTpl.find('.add-question-btn')
                .prop('disabled', false)
                .text('Добавить')
                .removeClass('btn-secondary')
                .addClass('btn-primary');

            fillQuestionTypeSelect($newTpl.find('select.question-type-select'), '');

            $tplRow.after($newTpl);
        });


        $('#deleteQuestion').on('show.bs.modal', function (e) {
            const qid = $(e.relatedTarget).closest('tr').data('question-id') || 0;
            $('#questionId').val(qid);
        });

        $(document).on('click', '#confirmDeleteQuestion', function () {
            const id = toInt($('#questionId').val());
            if (id) $(`tr[data-question-id="${id}"]`).remove();
            $('#deleteQuestion').modal('hide');
        });

        $('[id^="editOffer-"]').each(function() {
            const offerId = toInt(this.dataset.offerId);
            refreshFilterOptions(offerId);
        });

        $(document).on('click', '.btn-download-avito', function () {
            const vacId = parseInt($('#vacancyId').val() || '0', 10);
            if (!vacId) {
                showError('Сначала сохраните РК, чтобы выгрузить ID.');
                return;
            }

            fetch(`/rk/avito_ids?vacancyId=${vacId}`)
                .then(r => {
                    if (!r.ok) return r.text().then(t => { throw new Error(t || 'Ошибка выгрузки'); });
                    return r.blob();
                })
                .then(blob => {
                    const url = window.URL.createObjectURL(blob);
                    const a   = document.createElement('a');
                    a.href     = url;
                    a.download = `avito_ids_${vacId}.xlsx`;
                    document.body.appendChild(a);
                    a.click();
                    a.remove();
                    window.URL.revokeObjectURL(url);
                })
                .catch(err => showError(err.message));
        });

        $(document).on('click', '.manual-avito-btn', function () {
            currentManualCtx = {
                fileInput : '#' + $(this).data('file-input'),
                fileLabel : '#' + $(this).data('file-label'),
                infoDiv   : '#' + $(this).data('info-div'),
                button    : this
            };
            $('#manualAvitoIds').val(manualAvitoIds);
        });

        $(document).on('click', '#confirmManualAvito', function () {
            const raw = $('#manualAvitoIds').val().trim();
            if (!raw) {
                alert('Поле не должно быть пустым');
                return;
            }

            const vacId = parseInt($('#vacancyId').val() || '0', 10);
            const fd = new FormData();
            fd.append('inputType', 'manual');
            fd.append('manualAvitoIds', raw);

            fetch(`/rk/check_avito_ids?vacancyId=${vacId}`, {
                method: 'POST',
                body: fd
            })
                .then(r => {
                    if (!r.ok) {
                        return r.text().then(t => { throw new Error(t || 'Ошибка проверки Avito-ID'); });
                    }
                    return r.json();
                })
                .then(({ new: cntNew, count, total }) => {
                    const qty = cntNew ?? count ?? total ?? 0;

                    manualAvitoIds = raw;

                    $(currentManualCtx.infoDiv)
                        .removeClass('d-none')
                        .text(`Загружено ${qty} новых ID (ручной ввод)`);

                    $(currentManualCtx.fileInput).prop('disabled', true).val('');
                    $(currentManualCtx.fileLabel).addClass('disabled');

                    $('#dynamicErrorAlert').remove();
                    $('#manualAvitoModal').modal('hide');
                })
                .catch(err => {
                    $('#manualAvitoModal').modal('hide');
                    showError(err.message);
                });
        });
    }

    let allFilterOptions = "";
    (function collectAllFilterOptions() {
        $('tr[id^="newFilterRow-"] select option').each(function() {
            const val = $(this).attr('value');
            if (val) {
                allFilterOptions += this.outerHTML || this.cloneNode(true).outerHTML;
            }
        });
    })();

    let exportTypesList = [];
    let reklsList          = [];
    let buyingTypesList    = [];
    let professionsList    = [];
    let providersList      = [];
    let trafficSourcesList = [];

    async function loadDictionaries() {
        try {
            const [
                exportTypesData,
                questionTypesData,
                reklsData,
                buyingTypesData,
                professionsData,
                providersData,
                trafficSourcesData
            ] = await Promise.all([
                fetchJSON('/rk/export_types'),
                fetchJSON('/rk/question_types'),
                fetchJSON('/rk/rekls'),
                fetchJSON('/rk/buying_types'),
                fetchJSON('/rk/professions'),
                fetchJSON('/rk/providers'),
                fetchJSON('/rk/traffic_sources')
            ]);

            const normalize = arr =>
                Array.isArray(arr)
                    ? arr.map(o => ({ Id: o.Id ?? o.id, Name: o.Name ?? o.name }))
                    : [];

            exportTypesList   = normalize(exportTypesData);
            reklsList         = normalize(reklsData);
            buyingTypesList   = normalize(buyingTypesData);
            professionsList   = normalize(professionsData);
            providersList     = normalize(providersData);
            trafficSourcesList= normalize(trafficSourcesData);
            questionTypesList = normalize(questionTypesData);

            $('.export-type-select').each(function () {
                const $sel = $(this);
                const current = $sel.attr('data-current') || "";
                $sel.empty();
                $sel.append('<option value="">Выбрать тип экспорта</option>');
                exportTypesList.forEach(function (et) {
                    $sel.append(`<option value="${et.Id}">${et.Name}</option>`);
                });
                if (current !== "") {
                    $sel.val(current);
                }
            });

            $('select.question-type-select').each(function () {
                const $sel = $(this);
                fillQuestionTypeSelect($sel, $sel.data('current') || $sel.val());
            });

            renderVacancyDicts();
        } catch (e) {
            console.error('Не удалось загрузить справочники clients/export_types', e);
        }
    }

    function fillInitialFilterNames() {
        $('[data-filter-id]').each(function () {
            const fid  = parseInt(this.dataset.filterId, 10);
            const meta = allFilters.find(f => f.id === fid);
            if (!meta) return;
            $(this).find('.filter-name').text(meta.name  || '—');
            $(this).find('.filter-val' ).text(meta.value || '—');
        });
    }

    function fillQuestionTypeSelect($sel, value) {
        $sel.empty()
            .append('<option value="">Выбрать тип…</option>');
        questionTypesList.forEach(t =>
            $sel.append(`<option value="${t.Id}">${t.Name}</option>`));
        if (value !== undefined && value !== null) {
            $sel.val(String(value));
        }
        enableSelect2($sel);
    }

    loadDictionaries().then(async () => {
        await loadAllFilters();
        fillInitialFilterNames();
        initHandlers();

        saveState.vacancy = gatherVacancy();

        saveState.offers = $('[id^="editOffer-"]').map(function() {
            const offerId = toInt(this.dataset.offerId);
            return collectOfferFromDOM(offerId);
        }).get();

        saveState.offers.forEach(o => {
            rebuildFilterSelect(o.id);
        });

        if (!saveState.vacancy.openAiSupport) {
            saveState.questionnaire = gatherQuestionnaire();
        }
    });

    showToastAfterSaving('РК: ' +  $('#vacancyId').val() + ' успешно сохранилась')
});