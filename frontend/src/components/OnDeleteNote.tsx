import NoteService from "../service/NoteService";
import TagService from "../service/TagService";
import Tag from "../types/Tag";

export function OnDeleteNote(id: string, tags: Tag[]) {
    // try {
    //     TagService.detach(id, tag.id).then(
    //         (response) => { console.log("detached") },
    //         (error) => { console.log("error") },
    //     )
    // }
    tags.forEach(function (tag) {
        TagService.detach(id, tag.id).then(
            (response) => { console.log("detached") },
            (error) => { console.log("error") },
        )
    })
    try {
        NoteService.remove(id)
    } catch (error) {
        console.log(error)
    }
}
