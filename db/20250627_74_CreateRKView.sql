CREATE OR REPLACE VIEW rk_view AS
SELECT
    va.id,
    concat_ws(
            '|'::text,
            va.id::text,
            va.name,
            COALESCE(rdr.name,  ''::character varying),
            COALESCE(rdp2.name, ''::character varying),
            COALESCE(rdp.name,  ''::character varying),
            COALESCE(et.short_name, ''::character varying),
            COALESCE(rdbt.name, ''::character varying),
            COALESCE(rdts.name, ''::character varying),
            'чатбот',
            COALESCE(va.cpa_ex_id::text, ''::text),
            COALESCE(va.other_specification, ''::character varying)
    )::character varying                                 AS name,
    va.load_to_zp,
    va."clientId",
    va.quest_id,
    va.hello_text,
    va.bye_text,
    (
        SELECT string_agg(concat(v.id, ' : ', v.name), ', ' ORDER BY v.name)
        FROM vacancy_2_vacancy vv
                 JOIN vacancy v ON v.id = vv.vacancy_id
        WHERE vv.vacancy_avito_id = va.id
    )                                                    AS offers,
    q.name                                               AS questname,
    va.post_processing_messages,
    va.dialog_life_time_in_minutes,
    va.followup_message_interval_in_minutes
FROM vacancy_avito va
         JOIN export_type              et   ON et.id = va.export_type_id
         LEFT JOIN rk_dict_buying_type rdbt ON rdbt.id = va.buying_type_id
         LEFT JOIN rk_dict_provider    rdp  ON rdp.id = va.provider_id
         LEFT JOIN rk_dict_prof        rdp2 ON rdp2.id = va.prof_id
         LEFT JOIN rk_dict_rekl        rdr  ON rdr.id = va.rekl_id
         LEFT JOIN rk_dict_traffic_source rdts ON rdts.id = va.traffic_source_id
         LEFT JOIN questionnaire       q    ON q.id = va.quest_id
WHERE va.openai_support = false;