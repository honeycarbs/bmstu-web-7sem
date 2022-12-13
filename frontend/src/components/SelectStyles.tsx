export const CustomSelectStyle = {
    control: (baseStyles: any, state: any) => ({
        ...baseStyles,
        // borderColor: state.isFocused ? 'grey' : 'grey',
        borderRadius: '16px',
        backgroundColor: 'rgba(250, 250, 249, 255)',
        border: "2px solid rgba(205, 191, 180, 0.2)!important",
        boxShadow: 'none',
        minHeight: '45px',
        // minHeight: '45px',
        '&:focus-within': {
            color: '#BBAC9E!important',
        },
    }),
    multiValue: (baseStyles: any) => ({
        ...baseStyles,
        backgroundColor: '#BBAC9F!important',
        borderRadius: '16px',
        padding: '2px 5px'
    }),

    multiValueLabel: (baseStyles: any) => ({
        ...baseStyles,
        backgroundColor: '#BBAC9F!important',
        textAlign: "center",
        padding: '2px 5px',
        color: '#FFF!important',
        borderRadius: '16px',
    }),

    multiValueRemove: (baseStyles: any) => ({
        ...baseStyles,
        borderRadius: '16px',
        color: 'white',
    }),

    menu: (baseStyles: any) => ({
        ...baseStyles,
        borderRadius: '10px',
        opacity: '0.9',
        backgroundColor: 'rgba(250, 250, 249, 255)',
        border: "2px solid rgba(205, 191, 180, 0.2)!important",
        boxShadow: 'none',
    }),

    option: (baseStyles: any, state: { isSelected: any; }) => ({
        ...baseStyles,
        borderRadius: '16px',
        border: "0px",
        boxShadow: 'none',
        backgroundColor: "rgba(250, 250, 249, 255)",
        '&:hover': {
            backgroundColor: 'rgba(205, 191, 180, 0.4)',
        },
    }),
}