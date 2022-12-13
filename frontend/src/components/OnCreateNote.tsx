import NoteService from "../service/NoteService"
import TagService from "../service/TagService"
import CreateNoteDTO from "../types/dto/CreateNoteDTO"
import Note from "../types/Note"

export function OnCreateNote({ header, body, tags }: Note) {
    var dto: CreateNoteDTO = {
        header: header,
        body: body,
        color: "47B5FF"
    }

    NoteService.create(dto).then(
        (response) => {
            return response
        },
        (error) => {
            console.log(error)
        }
    ).then(
        (response) => {
            tags.map(tag => {
                TagService.create(String(response), tag).then(
                    (response) => {
                        return response
                    },
                    (error) => {
                        console.log("sraka")
                    }
                )
            })
        },
        (error) => {
            console.log("sraka")
        }
    )
}
