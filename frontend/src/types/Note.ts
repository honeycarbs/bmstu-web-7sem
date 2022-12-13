import Tag from "./Tag"

export default interface Note {
    id?: any | null,
    header: string,
    body: string,
    tags: Tag[]
    // color: string,
    // edited: string,
}

export type NoteRenderData = {
    id: number
    header: string
    tags: Tag[]
}