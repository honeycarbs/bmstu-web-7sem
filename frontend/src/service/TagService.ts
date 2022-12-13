import http from "../http-common";
import Tag from "../types/Tag";
import authHeader from "../authHeader";

const getAll = async () => {
    console.log("tags fetched")
    const response = http.get<Array<Tag>>("api/v1/tags", { headers: authHeader() });

    return (await response).data
};

const get = async (id: any) => {
    const response = http.get<Tag>(`api/v1/tags/${id}`, { headers: authHeader() });

    return (await response).data
};

const create = async (noteURL: string, data: Tag) => {
    try {
        const response = await http.post<any>(`${noteURL}/tags`, data, { headers: authHeader() });
        return response.data
    } catch (error) {
        console.log(error)
        return null
    }
};

const update = async (id: any, data: Tag) => {
    const response = http.patch<any>(`api/v1/tags/${id}`, data, { headers: authHeader() });

    return (await response).data
};

const detach = async (noteID: any, tagID: any) => {
    const response = http.delete<any>(`api/v1/notes/${noteID}/tags/${tagID}`, { headers: authHeader() });

    return (await response).data
};

const TagService = {
    getAll,
    get,
    create,
    update,
    detach,
};

export default TagService;