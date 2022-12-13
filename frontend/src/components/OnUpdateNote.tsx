import { useEffect, useState } from "react"
import NoteService from "../service/NoteService"
import TagService from "../service/TagService"
import Note from "../types/Note"
import Tag from "../types/Tag"

export interface UpdateNoteProps {
    note: Note,
    prevTags: Tag[]
}

// function isPresent

export function OnUpdateNote(note: Note, prevTags: Tag[]) {

    prevTags.forEach(function (tag) {
        const present = note.tags.find(t => t.id === tag.id)
        if (present == null) {
            TagService.detach(note.id, tag.id).then(
                (response) => { console.log("detached") },
                (error) => { console.log("error") },
            )
        }
    })

    note.tags.forEach(function (tag) {
        if (tag.id == 0) {
            console.log("i am a new tag: ", tag.label)
            TagService.create(`api/v1/notes/${note.id}`, tag)
        } else {
            const present = prevTags.find(t => t.id === tag.id)
            if (present == null) { 
                console.log("i am new to this note but i exist: ", tag.label)
                TagService.create(`api/v1/notes/${note.id}`, tag)
            }
        }
    })

    NoteService.update(note.id, note).then(
        (response) => { console.log("updated") },
        (error) => { console.log("error") },
    )
}
