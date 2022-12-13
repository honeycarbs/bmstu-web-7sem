import http from "../http-common";
import CreateNoteDTO from "../types/dto/CreateNoteDTO";
import Note from "../types/Note";
import authHeader from "../authHeader";
import { useState } from "react";

const getAll = async (searchUrl: string) => {
    console.log("notes fetched")

    const response = http.get<Array<Note>>("api/v1/notes" + searchUrl, { headers: authHeader() });

    return (await response).data
};

const get = async (url: string) => {
    try {
        const response = await http.get<Note>(url, { headers: authHeader() });

        if (response.status != 200) throw response.statusText;

        return response.data
    } catch (e) {
        console.log("here");
        throw e;
    }
    // const response = await http.get<Note>(url, { headers: authHeader() });

    // return response.data
};


const create = async (data: CreateNoteDTO) => {
    const response = http.post<string>("api/v1/notes", data, { headers: authHeader() });

    return (await response).data
};

const update = async (id: any, data: Note) => {
    const response = http.patch<any>(`api/v1/notes/${id}`, data, { headers: authHeader() });

    return (await response).data
};

const remove = async (id: any) => {
    const response = http.delete<any>(`api/v1/notes/${id}`, { headers: authHeader() });

    return (await response).data
};

const NoteService = {
    getAll,
    get,
    create,
    update,
    remove,
};

export default NoteService;