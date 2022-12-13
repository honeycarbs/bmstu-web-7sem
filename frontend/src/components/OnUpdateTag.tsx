import TagService from "../service/TagService"
import Tag from "../types/Tag"

export function OnUpdateTag(id: string, tag: Tag) {
    // console.log("i am updating tags")
    try {
        TagService.update(id, tag)
    } catch (error) {
        console.log(error)
    }
}